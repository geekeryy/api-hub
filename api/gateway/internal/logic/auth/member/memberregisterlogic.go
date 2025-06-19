package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册
func NewMemberRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRegisterLogic {
	return &MemberRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberRegisterLogic) MemberRegister(req *types.MemberRegisterReq) (resp *types.MemberRegisterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
