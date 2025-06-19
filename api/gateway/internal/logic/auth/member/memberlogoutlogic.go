package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登出
func NewMemberLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberLogoutLogic {
	return &MemberLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberLogoutLogic) MemberLogout(req *types.MemberLogoutReq) error {
	// todo: add your logic here and delete this line

	return nil
}
