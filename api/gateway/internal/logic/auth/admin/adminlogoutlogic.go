package admin

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登出
func NewAdminLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLogoutLogic {
	return &AdminLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLogoutLogic) AdminLogout(req *types.AdminLogoutReq) error {
	// todo: add your logic here and delete this line

	return nil
}
