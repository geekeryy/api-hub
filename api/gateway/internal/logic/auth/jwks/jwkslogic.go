package jwks

import (
	"context"
	"net/http"

	"github.com/MicahParks/jwkset"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/zeromicro/go-zero/core/logx"
)

type JWKSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取公钥
func NewJWKSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JWKSLogic {
	return &JWKSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JWKSLogic) JWKS(w http.ResponseWriter) error {
	jwksets := jwkset.NewMemoryStorage()
	jwksList, err := l.svcCtx.JwksModel.FindAll(l.ctx)
	if err != nil {
		l.Errorf("Failed to find all jwks publics. Error: %s", err)
		return err
	}
	if len(jwksList) > 0 {
		for _, record := range jwksList {
			publicKey, err := xstrings.AesCbcDecryptBase64(record.PublicKey, "public_key_secre", nil)
			if err != nil {
				l.Errorf("Failed to decrypt public key. Error: %s", err)
				return err
			}
			jwks.AddKey(l.ctx, record.Kid, []byte(publicKey), jwksets)
		}
	}
	rawJWKS, err := jwksets.JSONPublic(l.ctx)
	if err != nil {
		l.Errorf("Failed to get the server's JWKS. Error: %s", err)
		return err
	}
	w.Write(rawJWKS)
	return nil
}
