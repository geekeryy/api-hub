package svc

import (
	"context"
	"log"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/api/gateway/internal/config"
	"github.com/geekeryy/api-hub/api/gateway/internal/middleware"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xgorm"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/time/rate"
)

type ServiceContext struct {
	Config            config.Config
	ContextMiddleware rest.Middleware
	JwtMiddleware     rest.Middleware
	Validator         *validate.Validate
	Jwkset            *jwkset.MemoryJWKSet
	JwksModel         authmodel.JwksModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	pg, err := xgorm.ConnectPg(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to connect to database.\nError: %s", err)
	}
	jwksModel := authmodel.NewJwksModel(pg)

	jwksets := getJwkset(jwksModel)

	kfunc, err := jwks.InitKeyfunc(context.Background(), c.Jwks.ServerURL, keyfunc.Override{
		RefreshInterval:   time.Duration(c.Jwks.RefreshInterval) * time.Second,
		RateLimitWaitMax:  time.Duration(c.Jwks.RefreshInterval/2) * time.Second,
		RefreshUnknownKID: rate.NewLimiter(rate.Every(1*time.Minute), 2),
	})
	if err != nil {
		log.Fatalf("Failed to init keyfunc.\nError: %s", err)
	}

	svc := &ServiceContext{
		Config:            c,
		ContextMiddleware: middleware.NewContextMiddleware().Handle,
		JwtMiddleware:     middleware.NewJwtMiddleware(kfunc).Handle,
		Validator:         validate.New(nil, []string{"zh", "en"}),
		Jwkset:            jwksets,
		JwksModel:         jwksModel,
	}
	return svc
}

func getJwkset(jwksModel authmodel.JwksModel) *jwkset.MemoryJWKSet {
	jwksets := jwkset.NewMemoryStorage()
	jwksList, err := jwksModel.FindAll()
	if err != nil {
		log.Fatalf("Failed to find all jwks publics.\nError: %s", err)
	}
	if len(jwksList) > 0 {
		for _, record := range jwksList {
			publicKey, err := xstrings.AesCbcDecryptBase64(record.PublicKey, "public_key_secre", nil)
			if err != nil {
				log.Fatalf("Failed to decrypt public key.\nError: %s", err)
			}
			jwks.AddKey(context.Background(), record.Kid, []byte(publicKey), jwksets)
		}
	}
	return jwksets
}
