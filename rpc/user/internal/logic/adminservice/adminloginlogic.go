package adminservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/user"

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

func (l *AdminLoginLogic) AdminLogin(in *user.AdminLoginReq) (*user.AdminLoginResp, error) {
	// todo: add your logic here and delete this line

	return &user.AdminLoginResp{}, nil
}
