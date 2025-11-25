// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/user/client/authservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberRefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新Token
func NewMemberRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRefreshTokenLogic {
	return &MemberRefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberRefreshTokenLogic) MemberRefreshToken(req *types.MemberRefreshTokenReq) (resp *types.MemberRefreshTokenResp, err error) {
	// 验证refresh token
	claims, err := jwks.ValidateToken(req.RefreshToken, l.svcCtx.Kfunc)
	if err != nil {
		l.Errorf("Failed to validate token. Error: %s", err)
		return nil, xerror.UnauthorizedErr
	}
	memberId, err := claims.GetSubject()
	if err != nil || len(memberId) == 0 {
		l.Errorf("validate claims err:%v %+v", err, claims)
		return nil, xerror.UnauthorizedErr
	}

	_, err = l.svcCtx.MemberService.MemberRefreshToken(l.ctx, &authservice.MemberRefreshTokenReq{
		RefreshToken: req.RefreshToken,
		MemberUuid:   memberId,
	})
	if err != nil {
		return nil, err
	}

	return
}
