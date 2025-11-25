// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/rpc/user/client/memberservice"

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

func (l *MemberRegisterLogic) MemberRegister(req *types.MemberRegisterReq) error {
	_, err := l.svcCtx.MemberService.MemberRegister(l.ctx, &memberservice.MemberRegisterReq{
		IdentityType: req.IdentityType,
		Identifier:   req.Identifier,
		Credential:   req.Credential,
		Code:         req.Code,
	})
	if err != nil {
		return err
	}
	return nil
}
