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
	_, err := l.svcCtx.MemberService.MemberForgetPassword(l.ctx, &memberservice.MemberForgetPasswordReq{
		IdentityType: req.IdentityType,
		Identifier:   req.Identifier,
		Code:         req.Code,
		Password:     req.Password,
	})
	if err != nil {
		return err
	}
	return nil
}
