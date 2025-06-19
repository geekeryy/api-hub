package svc

import (
	"context"
	"log"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/api/gateway/internal/config"
	"github.com/geekeryy/api-hub/api/gateway/internal/middleware"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xgorm"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/time/rate"
)

type ServiceContext struct {
	Config            config.Config
	ContextMiddleware rest.Middleware
	JwtMiddleware     rest.Middleware
	Validator         *validate.Validate
	JwksModel         authmodel.JwksModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	pg, err := xgorm.ConnectPg(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %s", err)
	}

	kfunc, err := jwks.InitKeyfunc(context.Background(), c.Jwks.ServerURL, keyfunc.Override{
		RefreshInterval:   time.Duration(c.Jwks.RefreshInterval) * time.Second,
		RateLimitWaitMax:  time.Duration(c.Jwks.RefreshInterval/2) * time.Second,
		RefreshUnknownKID: rate.NewLimiter(rate.Every(1*time.Minute), 2),
	})
	if err != nil {
		log.Fatalf("Failed to init keyfunc. Error: %s", err)
	}

	svc := &ServiceContext{
		Config:            c,
		ContextMiddleware: middleware.NewContextMiddleware().Handle,
		JwtMiddleware:     middleware.NewJwtMiddleware(kfunc).Handle,
		JwksModel: authmodel.NewJwksModel(pg),
		Validator: validate.New([]validate.ValidatorFn{
			validator.ChineseNameValidator,
		}, []string{"zh", "en"}),
	}
	return svc
}
