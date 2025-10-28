package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/rpc/auth/client/authservice"

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

// TODO: 高并发情况下，会导致重复注册，需要优化
func (l *MemberRegisterLogic) MemberRegister(req *types.MemberRegisterReq) error {
	_, err := l.svcCtx.AuthService.MemberRegister(l.ctx, &authservice.MemberRegisterReq{
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
