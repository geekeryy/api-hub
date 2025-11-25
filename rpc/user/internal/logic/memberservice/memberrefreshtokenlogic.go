package memberservicelogic

import (
	"context"
	"time"

	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/model"
	"github.com/geekeryy/api-hub/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberRefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRefreshTokenLogic {
	return &MemberRefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MemberRefreshTokenLogic) MemberRefreshToken(in *user.MemberRefreshTokenReq) (*user.MemberRefreshTokenResp, error) {
	refreshTokenHash, err := xstrings.AesCbcEncryptBase64(in.RefreshToken, l.svcCtx.Config.Secret.RefreshToken, nil)
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
		l.Errorf("Refresh token is disabled. RefreshToken: %s", in.RefreshToken)
		return nil, xerror.UnauthorizedErr
	}
	if refreshToken.MemberId != in.MemberUuid {
		l.Errorf("Refresh token member id not match. RefreshToken: %s", in.RefreshToken)
		return nil, xerror.UnauthorizedErr
	}

	// 生成新的token
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
	token, exp, err := jwks.GenerateToken(jwksRecord.Kid, refreshToken.MemberId, privateKey, int64(l.svcCtx.Config.Jwt.AccessExpire), nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}

	_, err = l.svcCtx.TokenRefreshRecordModel.Insert(l.ctx, nil, &model.TokenRefreshRecord{
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

	return &user.MemberRefreshTokenResp{
		Token: token,
	}, nil
}
