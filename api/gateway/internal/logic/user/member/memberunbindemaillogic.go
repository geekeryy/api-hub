package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberUnbindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 解绑邮箱
func NewMemberUnbindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberUnbindEmailLogic {
	return &MemberUnbindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberUnbindEmailLogic) MemberUnbindEmail(req *types.MemberUnbindEmailReq) error {
	// todo: add your logic here and delete this line

	return nil
}
