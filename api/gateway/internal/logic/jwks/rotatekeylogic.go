package jwks

import (
	"context"
	"log"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/rpc/model/jwksmodel"

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

func (l *RotateKeyLogic) RotateKey(req *types.RotateKeyReq) error {
	pub, priv, err := jwks.RotateKey(l.ctx, l.svcCtx.Jwkset)
	if err != nil {
		log.Fatalf("Failed to rotate key.\nError: %s", err)
	}
	l.svcCtx.RWMKey.Lock()
	defer l.svcCtx.RWMKey.Unlock()
	l.svcCtx.PrivateKey = priv
	l.svcCtx.PublicKey = pub

	encryptPub, err := xstrings.AesCbcEncryptBase64(string(pub), "public_key_secre", nil)
	if err != nil {
		log.Fatalf("Failed to encrypt public key.\nError: %s", err)
	}
	encryptPriv, err := xstrings.AesCbcEncryptBase64(string(priv), "private_key_secr", nil)
	if err != nil {
		log.Fatalf("Failed to encrypt private key.\nError: %s", err)
	}
	if err := l.svcCtx.JwksPublicModel.Insert(l.ctx, nil, &jwksmodel.JwksPublic{
		PublicKey:  encryptPub,
		PrivateKey: encryptPriv,
	}); err != nil {
		log.Fatalf("Failed to insert jwks public.\nError: %s", err)
	}

	return nil
}
