package member

import (
	"context"
	"time"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"

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
	refreshTokenHash, err := xstrings.AesCbcEncryptBase64(req.RefreshToken, l.svcCtx.Config.Secret.RefreshToken, nil)
	if err != nil {
		l.Errorf("Failed to encrypt refresh token. Error: %s", err)
		return nil, err
	}
	refreshToken, err := l.svcCtx.RefreshTokenModel.FindOneByRefreshTokenHash(l.ctx, refreshTokenHash)
	if err != nil {
		l.Errorf("Failed to find refresh token. Error: %s", err)
		return nil, xerror.InternalServerErr
	}
	if refreshToken.Status == consts.RefreshTokenStatusDisabled || refreshToken.ExpiredAt.Before(time.Now()) {
		l.Errorf("Refresh token is disabled. RefreshToken: %s", req.RefreshToken)
		return nil, xerror.UnauthorizedErr
	}
	if refreshToken.MemberId != memberId {
		l.Errorf("Refresh token member id not match. RefreshToken: %s", req.RefreshToken)
		return nil, xerror.UnauthorizedErr
	}

	// 生成新的token
	jwksRecord, err := l.svcCtx.JwksModel.FindLatest(l.ctx)
	if err != nil {
		l.Errorf("Failed to find latest jwks public. Error: %s", err)
		return nil, xerror.InternalServerErr
	}
	privateKey, err := xstrings.AesCbcDecryptBase64(jwksRecord.PrivateKey, "private_key_secr", nil)
	if err != nil {
		l.Errorf("Failed to decrypt private key. Error: %s", err)
		return nil, err
	}
	token, exp, err := jwks.GenerateToken(jwksRecord.Kid, refreshToken.MemberId, privateKey, int64(l.svcCtx.Config.Auth.AccessExpire), nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}

	err = l.svcCtx.TokenRefreshRecordModel.Insert(l.ctx, nil, &authmodel.TokenRefreshRecord{
		RefreshTokenHash: refreshTokenHash,
		Token:            token,
		Kid:              jwksRecord.Kid,
		Ip:               xcontext.GetClientIp(l.ctx),
		ExpiredAt:        exp,
	})
	if err != nil {
		l.Errorf("Failed to insert token refresh record. Error: %s", err)
		return nil, xerror.InternalServerErr
	}

	resp = &types.MemberRefreshTokenResp{
		Token: token,
	}

	return
}
