package svc

import (
	"context"
	"fmt"
	"log"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/api/oms/internal/config"
	"github.com/geekeryy/api-hub/api/oms/internal/middleware"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/geekeryy/api-hub/rpc/user/client/memberservice"
	usermodel "github.com/geekeryy/api-hub/rpc/user/model"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/limiter"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xgrpc"
)

type ServiceContext struct {
	Config                  config.Config
	ContextMiddleware       rest.Middleware
	OmsOtpMiddleware        rest.Middleware
	OmsJwtMiddleware        rest.Middleware
	Validator               *validate.Validate
	JwksModel               usermodel.JwksModel
	TokenRefreshRecordModel usermodel.TokenRefreshRecordModel
	MemberInfoModel         usermodel.MemberInfoModel
	RefreshTokenModel       usermodel.RefreshTokenModel
	MemberService           memberservice.MemberService
	DB                      sqlx.SqlConn
	RedisClient             *redis.Client
	Kfunc                   keyfunc.Keyfunc
	CodeLimiter             *lru.Cache[string, *limiter.Limiter]
	GenerateTokenFunc       jwks.GenerateTokenFunc
	logx.Logger
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
	kfunc, generateTokenFunc, err := jwks.NewKeyfunc(context.Background())
	if err != nil {
		log.Fatalf("Failed to init keyfunc. Error: %s", err)
	}

	codeLimiter, err := lru.New[string, *limiter.Limiter](1000)
	if err != nil {
		log.Fatalf("Failed to init code limiter. Error: %s", err)
	}

	client := zrpc.MustNewClient(c.MemberService)

	logger := logx.WithContext(context.Background())

	monitorLazyClient := xgrpc.NewLazyClient(c.MonitorService, 60)

	svc := &ServiceContext{
		Config:            c,
		ContextMiddleware: middleware.NewContextMiddleware(monitorLazyClient, logger).Handle,
		OmsJwtMiddleware:  middleware.NewOmsJwtMiddleware(kfunc).Handle,
		Kfunc:             kfunc,
		GenerateTokenFunc: generateTokenFunc,
		Validator: validate.New([]validate.ValidatorFn{
			validator.ChineseNameValidator,
		}, []string{"zh", "en"}),
		RedisClient: redisClient,
		CodeLimiter: codeLimiter,

		JwksModel:               usermodel.NewJwksModel(mysqlClient),
		TokenRefreshRecordModel: usermodel.NewTokenRefreshRecordModel(mysqlClient),
		MemberInfoModel:         usermodel.NewMemberInfoModel(mysqlClient),
		RefreshTokenModel:       usermodel.NewRefreshTokenModel(mysqlClient),

		MemberService: memberservice.NewMemberService(client),
		Logger:        logger,
	}
	return svc
}

func (s *ServiceContext) Close() {
	// TODO graceful stop
}
