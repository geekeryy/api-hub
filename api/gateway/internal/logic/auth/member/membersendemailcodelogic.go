package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberSendEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送邮箱验证码
func NewMemberSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberSendEmailCodeLogic {
	return &MemberSendEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberSendEmailCodeLogic) MemberSendEmailCode(req *types.MemberSendEmailCodeReq) error {
	// todo: add your logic here and delete this line

	return nil
}
