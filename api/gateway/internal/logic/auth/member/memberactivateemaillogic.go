package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberActivateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 激活邮箱
func NewMemberActivateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberActivateEmailLogic {
	return &MemberActivateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberActivateEmailLogic) MemberActivateEmail(req *types.MemberActivateEmailReq) error {
	// todo: add your logic here and delete this line

	return nil
}
