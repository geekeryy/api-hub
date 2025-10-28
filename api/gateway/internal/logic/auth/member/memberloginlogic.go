package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/rpc/auth/client/authservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录
func NewMemberLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberLoginLogic {
	return &MemberLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TODO 多个账号绑定同一个第三方账号/手机号/邮箱，需要处理
func (l *MemberLoginLogic) MemberLogin(req *types.MemberLoginReq) (resp *types.MemberLoginResp, err error) {
	_, err = l.svcCtx.AuthService.MemberLogin(l.ctx, &authservice.MemberLoginReq{
		IdentityType: req.IdentityType,
		Identifier:   req.Identifier,
		Credential:   req.Credential,
	})
	if err != nil {
		return nil, err
	}
	return
}
