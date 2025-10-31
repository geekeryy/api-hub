package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/api/gateway/internal/config"
	"github.com/geekeryy/api-hub/api/gateway/internal/middleware"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/limiter"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xgrpc"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/geekeryy/api-hub/rpc/auth/client/authservice"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/geekeryy/api-hub/rpc/model/membermodel"
	"github.com/geekeryy/api-hub/rpc/monitor/client/monitorservice"
	"github.com/geekeryy/api-hub/rpc/user/client/memberservice"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/time/rate"
)

type ServiceContext struct {
	Config                  config.Config
	ContextMiddleware       rest.Middleware
	JwtMiddleware           rest.Middleware
	AdminJwtMiddleware      rest.Middleware
	OmsOtpMiddleware        rest.Middleware
	OmsJwtMiddleware        rest.Middleware
	Validator               *validate.Validate
	JwksModel               authmodel.JwksModel
	TokenRefreshRecordModel authmodel.TokenRefreshRecordModel
	MemberIdentityModel     membermodel.MemberIdentityModel
	MemberInfoModel         membermodel.MemberInfoModel
	RefreshTokenModel       authmodel.RefreshTokenModel
	MemberService           memberservice.MemberService
	AuthService             authservice.AuthService
	DB                      sqlx.SqlConn
	Kfunc                   keyfunc.Keyfunc
	CodeLimiter             *lru.Cache[string, *limiter.Limiter]
	RedisClient             *redis.Client
	MonitorService          monitorservice.MonitorService
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlClient, err := sqlx.NewConn(sqlx.SqlConf{
		DataSource: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Mysql.Username, c.Mysql.Password, c.Mysql.Host, c.Mysql.Dbname),
		DriverName: "mysql",
		Replicas:   nil,
		Policy:     "",
	})
	if err != nil {
		log.Fatalf("failed to open mysql: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	})
	if err != nil {
		log.Fatalf("failed to open redis: %v", err)
	}

	codeLimiter, err := lru.New[string, *limiter.Limiter](1000)
	if err != nil {
		log.Fatalf("Failed to init code limiter. Error: %s", err)
	}

	userClient := zrpc.MustNewClient(c.UserService)
	authClient := zrpc.MustNewClient(c.AuthService)
	monitorLazyClient := xgrpc.NewLazyClient(c.MonitorService, 3)

	kfunc, err := jwks.NewDefaultOverrideCtx(context.Background(), requestJWKSetFromGrpc(authClient), keyfunc.Override{
		RefreshInterval:   time.Duration(c.Jwks.RefreshInterval) * time.Second,
		RateLimitWaitMax:  time.Duration(c.Jwks.RefreshInterval/2) * time.Second,
		RefreshUnknownKID: rate.NewLimiter(rate.Every(1*time.Minute), 2),
	})
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}

	svc := &ServiceContext{
		Config:                  c,
		ContextMiddleware:       middleware.NewContextMiddleware(monitorLazyClient).Handle,
		JwtMiddleware:           middleware.NewJwtMiddleware(kfunc).Handle,
		AdminJwtMiddleware:      middleware.NewAdminJwtMiddleware(kfunc).Handle,
		Kfunc:                   kfunc,
		JwksModel:               authmodel.NewJwksModel(mysqlClient),
		TokenRefreshRecordModel: authmodel.NewTokenRefreshRecordModel(mysqlClient),
		MemberIdentityModel:     membermodel.NewMemberIdentityModel(mysqlClient),
		MemberInfoModel:         membermodel.NewMemberInfoModel(mysqlClient),
		RefreshTokenModel:       authmodel.NewRefreshTokenModel(mysqlClient),
		Validator: validate.New([]validate.ValidatorFn{
			validator.ChineseNameValidator,
		}, []string{"zh", "en"}),
		DB:             mysqlClient,
		RedisClient:    redisClient,
		CodeLimiter:    codeLimiter,
		MemberService:  memberservice.NewMemberService(userClient),
		AuthService:    authservice.NewAuthService(authClient),
		MonitorService: monitorservice.NewMonitorService(monitorLazyClient),
	}
	return svc
}

func (s *ServiceContext) Close() {

}

func requestJWKSetFromGrpc(conn zrpc.Client) func(ctx context.Context) (jwkset.JWKSMarshal, error) {
	return func(ctx context.Context) (jwkset.JWKSMarshal, error) {
		response, err := authservice.NewAuthService(conn).GetJwks(ctx, &authservice.GetJwksReq{})
		if err != nil {
			return jwkset.JWKSMarshal{}, err
		}
		var jwks jwkset.JWKSMarshal
		err = json.Unmarshal([]byte(response.Data), &jwks)
		if err != nil {
			return jwkset.JWKSMarshal{}, err
		}
		return jwks, nil
	}
}
