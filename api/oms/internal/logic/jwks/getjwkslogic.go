// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package jwks

import (
	"context"

	"github.com/MicahParks/jwkset"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
	"github.com/geekeryy/api-hub/api/oms/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJWKSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// jwks密钥列表
func NewGetJWKSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJWKSLogic {
	return &GetJWKSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJWKSLogic) GetJWKS() (resp *types.GetJWKSResp, err error) {
	jwksets := jwkset.NewMemoryStorage()
	jwksList, err := l.svcCtx.JwksModel.FindAll(l.ctx)
	if err != nil {
		l.Errorf("Failed to find all jwks publics. Error: %s", err)
		return nil, err
	}
	if len(jwksList) > 0 {
		for _, record := range jwksList {
			publicKey, err := xstrings.AesCbcDecryptBase64(record.PublicKey, l.svcCtx.Config.Secret.PublicKey, nil)
			if err != nil {
				l.Errorf("Failed to decrypt public key. Error: %s", err)
				return nil, err
			}
			jwks.AddKey(l.ctx, record.Kid, []byte(publicKey), jwksets)
		}
	}

	keys, err := jwksets.MarshalWithOptions(l.ctx, jwkset.JWKMarshalOptions{
		Private: true,
	}, jwkset.JWKValidateOptions{})
	if err != nil {
		l.Errorf("Failed to get the server's JWKS keys. Error: %s", err)
		return nil, err
	}
	resp = &types.GetJWKSResp{
		Data: make([]types.GetJWKSItem, 0),
	}
	for _, key := range keys.Keys {
		resp.Data = append(resp.Data, types.GetJWKSItem{
			Kid: key.KID,
			Kty: key.KTY.String(),
			Alg: key.ALG.String(),
			Crv: key.CRV.String(),
			X:   key.X,
		})
	}

	return resp, nil
}
