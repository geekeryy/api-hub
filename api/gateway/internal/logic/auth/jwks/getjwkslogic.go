package jwks

import (
	"context"
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/rpc/auth/client/authservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJWKSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取公钥
func NewGetJWKSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJWKSLogic {
	return &GetJWKSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJWKSLogic) GetJWKS(req *types.JWKSReq, w http.ResponseWriter) error {
	resp, err := l.svcCtx.AuthService.GetJwks(l.ctx, &authservice.GetJwksReq{
		Service: req.Service,
	})
	if err != nil {
		return err
	}
	w.Write([]byte(resp.Data))
	return nil
}
