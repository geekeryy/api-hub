package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberBindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定手机号
func NewMemberBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberBindPhoneLogic {
	return &MemberBindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberBindPhoneLogic) MemberBindPhone(req *types.MemberBindPhoneReq) error {
	// todo: add your logic here and delete this line

	return nil
}
