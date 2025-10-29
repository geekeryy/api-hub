package auth

import (
	"context"

	"github.com/geekeryy/api-hub/api/oms/internal/svc"
	"github.com/geekeryy/api-hub/api/oms/internal/types"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/pquerna/otp/totp"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	switch req.IdentityType {
	case consts.IdentityTypePhone:
		return nil, xerror.InvalidParameterErr
	case consts.IdentityTypeOtp:
		if !totp.Validate(req.Credential, l.svcCtx.Config.Oms.OtpSecret) {
			return nil, xerror.UnauthorizedErr
		}
		token, _, err := l.svcCtx.GenerateTokenFunc("admin", 10*60, nil)
		if err != nil {
			return nil, xerror.InternalServerErr
		}
		return &types.LoginResp{
			Token: token,
		}, nil

	default:
		return nil, xerror.InvalidParameterErr
	}

}
