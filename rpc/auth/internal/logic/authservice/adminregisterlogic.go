package authservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/rpc/auth/auth"
	"github.com/geekeryy/api-hub/rpc/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRegisterLogic {
	return &AdminRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminRegisterLogic) AdminRegister(in *auth.AdminRegisterReq) (*auth.AdminRegisterResp, error) {
	// todo: add your logic here and delete this line

	return &auth.AdminRegisterResp{}, nil
}
