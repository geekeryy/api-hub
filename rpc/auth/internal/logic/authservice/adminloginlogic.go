package authservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/rpc/auth/auth"
	"github.com/geekeryy/api-hub/rpc/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// B端用户授权
func (l *AdminLoginLogic) AdminLogin(in *auth.AdminLoginReq) (*auth.AdminLoginResp, error) {
	// todo: add your logic here and delete this line

	return &auth.AdminLoginResp{}, nil
}
