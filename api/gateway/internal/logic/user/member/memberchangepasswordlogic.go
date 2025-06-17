package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改密码
func NewMemberChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberChangePasswordLogic {
	return &MemberChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberChangePasswordLogic) MemberChangePassword(req *types.MemberChangePasswordReq) error {
	// todo: add your logic here and delete this line

	return nil
}
