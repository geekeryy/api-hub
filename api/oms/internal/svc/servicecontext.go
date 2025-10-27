package svc

import (
	"context"
	"log"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/api/oms/internal/config"
	"github.com/geekeryy/api-hub/api/oms/internal/middleware"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/geekeryy/api-hub/rpc/model/usermodel"
	"github.com/geekeryy/api-hub/rpc/user/client/memberservice"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/time/rate"
	"gorm.io/gorm"

	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/limiter"
	"github.com/geekeryy/api-hub/core/pgcache"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xgorm"
)

type ServiceContext struct {
	Config                  config.Config
	ContextMiddleware       rest.Middleware
	OmsOtpMiddleware        rest.Middleware
	OmsJwtMiddleware        rest.Middleware
	Validator               *validate.Validate
	JwksModel               authmodel.JwksModel
	TokenRefreshRecordModel authmodel.TokenRefreshRecordModel
	MemberIdentityModel     authmodel.MemberIdentityModel
	MemberInfoModel         usermodel.MemberInfoModel
	RefreshTokenModel       authmodel.RefreshTokenModel
	MemberService           memberservice.MemberService
	DB                      *gorm.DB
	Cache                   *pgcache.Cache
	Kfunc                   keyfunc.Keyfunc
	CodeLimiter             *lru.Cache[string, *limiter.Limiter]
}

func NewServiceContext(c config.Config) *ServiceContext {
	pg, err := xgorm.ConnectPg(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %s", err)
	}

	cache, err := pgcache.NewCache(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to init cache. Error: %s", err)
	}

	kfunc, err := jwks.InitKeyfunc(context.Background(), c.Jwks.ServerURL, keyfunc.Override{
		RefreshInterval:   time.Duration(c.Jwks.RefreshInterval) * time.Second,
		RateLimitWaitMax:  time.Duration(c.Jwks.RefreshInterval/2) * time.Second,
		RefreshUnknownKID: rate.NewLimiter(rate.Every(1*time.Minute), 2),
	})
	if err != nil {
		log.Fatalf("Failed to init keyfunc. Error: %s", err)
	}

	codeLimiter, err := lru.New[string, *limiter.Limiter](1000)
	if err != nil {
		log.Fatalf("Failed to init code limiter. Error: %s", err)
	}

	client := zrpc.MustNewClient(c.MemberService)

	svc := &ServiceContext{
		Config:                  c,
		ContextMiddleware:       middleware.NewContextMiddleware().Handle,
		OmsJwtMiddleware:        middleware.NewOmsJwtMiddleware(kfunc).Handle,
		Kfunc:                   kfunc,
		JwksModel:               authmodel.NewJwksModel(pg),
		TokenRefreshRecordModel: authmodel.NewTokenRefreshRecordModel(pg),
		MemberIdentityModel:     authmodel.NewMemberIdentityModel(pg),
		MemberInfoModel:         usermodel.NewMemberInfoModel(pg),
		RefreshTokenModel:       authmodel.NewRefreshTokenModel(pg),
		Validator: validate.New([]validate.ValidatorFn{
			validator.ChineseNameValidator,
		}, []string{"zh", "en"}),
		Cache:         cache,
		DB:            pg,
		CodeLimiter:   codeLimiter,
		MemberService: memberservice.NewMemberService(client),
	}
	return svc
}

func (s *ServiceContext) Close() {
	// TODO graceful stop
}
