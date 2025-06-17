package jwks

import (
	"context"
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
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
	rawJWKS, err := l.svcCtx.Jwkset.JSONPublic(l.ctx)
	if err != nil {
		l.Errorf("Failed to get the server's JWKS.\nError: %s", err)
		return err
	}
	l.Logger.Infof("JWKS: %s", string(rawJWKS))

	w.Write(rawJWKS)
	return nil
}
