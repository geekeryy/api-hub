package admin

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册
func NewAdminRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRegisterLogic {
	return &AdminRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminRegisterLogic) AdminRegister(req *types.AdminRegisterReq) (resp *types.AdminRegisterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
