package svc

import (
	"context"
	"log"
	"sync"

	"github.com/MicahParks/jwkset"
	"github.com/geekeryy/api-hub/api/gateway/internal/config"
	"github.com/geekeryy/api-hub/api/gateway/internal/middleware"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xgorm"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/rpc/model/jwksmodel"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config            config.Config
	ContextMiddleware rest.Middleware
	Validator         *validate.Validate
	Jwkset            *jwkset.MemoryJWKSet
	PrivateKey        []byte
	PublicKey         []byte
	RWMKey            sync.RWMutex
	JwksPublicModel   jwksmodel.JwksPublicModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	pg, err := xgorm.ConnectPg(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to connect to database.\nError: %s", err)
	}
	jwksPublicModel := jwksmodel.NewJwksPublicModel(pg)

	jwksets, pub, priv := getJwkset(jwksPublicModel)

	svc := &ServiceContext{
		Config:            c,
		ContextMiddleware: middleware.NewContextMiddleware().Handle,
		Validator:         validate.New(nil, []string{"zh", "en"}),
		Jwkset:            jwksets,
		PrivateKey:        priv,
		PublicKey:         pub,
		JwksPublicModel:   jwksPublicModel,
	}
	return svc
}

func getJwkset(jwksPublicModel jwksmodel.JwksPublicModel) (*jwkset.MemoryJWKSet, []byte, []byte) {
	var pub, priv []byte
	jwksets := jwkset.NewMemoryStorage()
	jwksPublics, err := jwksPublicModel.FindAll()
	if err != nil {
		log.Fatalf("Failed to find all jwks publics.\nError: %s", err)
	}
	if len(jwksPublics) > 0 {
		for i, jwksPublic := range jwksPublics {
			publicKey, err := xstrings.AesCbcDecryptBase64(jwksPublic.PublicKey, "public_key_secre", nil)
			if err != nil {
				log.Fatalf("Failed to decrypt public key.\nError: %s", err)
			}
			jwks.AddKey(context.Background(), []byte(publicKey), jwksets)

			if i == 0 {
				privateKey, err := xstrings.AesCbcDecryptBase64(jwksPublic.PrivateKey, "private_key_secr", nil)
				if err != nil {
					log.Fatalf("Failed to decrypt private key.\nError: %s", err)
				}
				pub = []byte(publicKey)
				priv = []byte(privateKey)
			}
		}
	} else {
		// TODO 加分布式锁，避免多个实例同时生成公钥
		pub, priv, err = jwks.RotateKey(context.Background(), jwksets)
		if err != nil {
			log.Fatalf("Failed to rotate key.\nError: %s", err)
		}
		encryptPub, err := xstrings.AesCbcEncryptBase64(string(pub), "public_key_secre", nil)
		if err != nil {
			log.Fatalf("Failed to encrypt public key.\nError: %s", err)
		}
		encryptPriv, err := xstrings.AesCbcEncryptBase64(string(priv), "private_key_secr", nil)
		if err != nil {
			log.Fatalf("Failed to encrypt private key.\nError: %s", err)
		}
		if err := jwksPublicModel.Insert(context.Background(), nil, &jwksmodel.JwksPublic{
			PublicKey:  encryptPub,
			PrivateKey: encryptPriv,
		}); err != nil {
			log.Fatalf("Failed to insert jwks public.\nError: %s", err)
		}
	}
	return jwksets, pub, priv
}
