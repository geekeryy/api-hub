package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberUnbindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 解绑手机号
func NewMemberUnbindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberUnbindPhoneLogic {
	return &MemberUnbindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberUnbindPhoneLogic) MemberUnbindPhone(req *types.MemberUnbindPhoneReq) error {
	// todo: add your logic here and delete this line

	return nil
}
