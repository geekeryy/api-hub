package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberForgetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 忘记密码
func NewMemberForgetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberForgetPasswordLogic {
	return &MemberForgetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberForgetPasswordLogic) MemberForgetPassword(req *types.MemberForgetPasswordReq) error {
	// todo: add your logic here and delete this line

	return nil
}
