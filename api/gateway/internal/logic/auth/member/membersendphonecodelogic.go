package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberSendPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送手机验证码
func NewMemberSendPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberSendPhoneCodeLogic {
	return &MemberSendPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberSendPhoneCodeLogic) MemberSendPhoneCode(req *types.MemberSendPhoneCodeReq) error {
	// todo: add your logic here and delete this line

	return nil
}
