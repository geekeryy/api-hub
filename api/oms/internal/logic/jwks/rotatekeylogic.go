package jwks

import (
	"context"

	"github.com/MicahParks/jwkset"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type RotateKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 轮换公钥
func NewRotateKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RotateKeyLogic {
	return &RotateKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RotateKeyLogic) RotateKey() error {
	jwksets := jwkset.NewMemoryStorage()
	kid := uuid.New().String()
	pub, priv, err := jwks.RotateKey(l.ctx, kid, jwksets)
	if err != nil {
		l.Errorf("Failed to rotate key. Error: %s", err)
		return err
	}
	encryptPub, err := xstrings.AesCbcEncryptBase64(string(pub), l.svcCtx.Config.Secret.PublicKey, nil)
	if err != nil {
		l.Errorf("Failed to encrypt public key. Error: %s", err)
		return err
	}
	encryptPriv, err := xstrings.AesCbcEncryptBase64(string(priv), l.svcCtx.Config.Secret.PrivateKey, nil)
	if err != nil {
		l.Errorf("Failed to encrypt private key. Error: %s", err)
		return err
	}
	if _, err := l.svcCtx.JwksModel.Insert(l.ctx, nil, &authmodel.Jwks{
		Kid:        kid,
		Service:    "",
		PublicKey:  encryptPub,
		PrivateKey: encryptPriv,
	}); err != nil {
		l.Errorf("Failed to insert jwks public. Error: %s", err)
		return err
	}
	return nil
}
