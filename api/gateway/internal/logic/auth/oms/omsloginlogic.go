package oms

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/pquerna/otp/totp"

	"github.com/zeromicro/go-zero/core/logx"
)

type OmsLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Oms登录
func NewOmsLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OmsLoginLogic {
	return &OmsLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OmsLoginLogic) OmsLogin(req *types.OmsLoginReq) (resp *types.OmsLoginResp, err error) {
	if req.Username != "admin" || !totp.Validate(req.Code, l.svcCtx.Config.Oms.OtpSecret) {
		return nil, xerror.ForbiddenErr
	}

	// 生成token
	jwksRecord, err := l.svcCtx.JwksModel.FindLatest(l.ctx)
	if err != nil {
		l.Errorf("Failed to find latest jwks public. Error: %s", err)
		return nil, xerror.InternalServerErr
	}
	privateKey, err := xstrings.AesCbcDecryptBase64(jwksRecord.PrivateKey, l.svcCtx.Config.Secret.PrivateKey, nil)
	if err != nil {
		l.Errorf("Failed to decrypt private key. Error: %s", err)
		return nil, err
	}
	token, _, err := jwks.GenerateToken(jwksRecord.Kid, "oms", privateKey, int64(l.svcCtx.Config.Auth.AccessExpire), nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}

	return &types.OmsLoginResp{
		Token: token,
	}, nil
}
