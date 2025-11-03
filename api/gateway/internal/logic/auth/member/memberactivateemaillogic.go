package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/auth/client/authservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberActivateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 激活邮箱
func NewMemberActivateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberActivateEmailLogic {
	return &MemberActivateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberActivateEmailLogic) MemberActivateEmail(req *types.MemberActivateEmailReq) error {
	claims, err := jwks.ValidateToken(req.Token, l.svcCtx.Kfunc)
	if err != nil {
		return xerror.New(err, xerror.InvalidParameterErr)
	}
	memberUUID, err := claims.GetSubject()
	if err != nil {
		return xerror.New(err, xerror.InvalidParameterErr)
	}
	email, err := jwks.MapClaimsParseString(claims, "email")
	if err != nil {
		return xerror.New(err, xerror.InvalidParameterErr)
	}

	_, err = l.svcCtx.AuthService.MemberActivateEmail(l.ctx, &authservice.MemberActivateEmailReq{
		MemberUuid: memberUUID,
		Email:      email,
	})
	if err != nil {
		return err
	}

	return nil
}
