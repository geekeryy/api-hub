package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"

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

func (l *MemberLoginLogic) MemberLogin(req *types.MemberLoginReq) (resp *types.MemberLoginResp, err error) {
	jwksRecord, err := l.svcCtx.JwksModel.FindLatest()
	if err != nil {
		l.Errorf("Failed to find latest jwks public. Error: %s", err)
		return nil, err
	}
	privateKey, err := xstrings.AesCbcDecryptBase64(jwksRecord.PrivateKey, "private_key_secre", nil)
	if err != nil {
		l.Errorf("Failed to decrypt private key. Error: %s", err)
		return nil, err
	}
	token, err := jwks.GenerateToken(jwksRecord.Kid, "1234567890", string(privateKey), l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}
	resp.Token = token
	resp.RefreshToken = token

	return
}
