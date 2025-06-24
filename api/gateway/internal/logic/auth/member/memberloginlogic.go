package member

import (
	"context"
	"fmt"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/facebook"
	"github.com/geekeryy/api-hub/core/google"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/geekeryy/api-hub/rpc/model/usermodel"
	"github.com/google/uuid"
	"gorm.io/gorm"

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
	var memberID string
	var thirdPartyId string
	var memberInfo *usermodel.MemberInfo
	switch req.IdentityType {
	case consts.IdentityTypePhone:
		cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("phone_code_%s", req.Identifier))
		if err != nil {
			l.Errorf("Failed to get cache. Error: %s", err)
			return nil, xerror.InternalServerErr
		}
		if req.Credential != cacheValue {
			l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
			return nil, xerror.UnauthorizedErr
		}
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, req.IdentityType, req.Identifier)
		if err != nil {
			l.Errorf("Failed to find member identity. Error: %s", err)
			return nil, xerror.InternalServerErr
		}

		if len(memberIdentities) == 0 {
			l.Errorf("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
			return nil, xerror.NotFoundErr
		}
		if err := l.svcCtx.Cache.Delete(fmt.Sprintf("phone_code_%s", req.Identifier)); err != nil {
			l.Errorf("Failed to delete cache. Error: %s", err)
		}
		memberID = memberIdentities[0].MemberId
	case consts.IdentityTypeEmail:
		cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("email_code_%s", req.Identifier))
		if err != nil {
			l.Errorf("Failed to get cache. Error: %s", err)
			return nil, xerror.InternalServerErr
		}
		if req.Credential != cacheValue {
			l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
			return nil, xerror.UnauthorizedErr
		}
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, req.IdentityType, req.Identifier)
		if err != nil {
			l.Errorf("Failed to find member identity. Error: %s", err)
			return nil, xerror.InternalServerErr
		}

		if len(memberIdentities) == 0 {
			l.Errorf("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
			return nil, xerror.NotFoundErr
		}
		if err := l.svcCtx.Cache.Delete(fmt.Sprintf("email_code_%s", req.Identifier)); err != nil {
			l.Errorf("Failed to delete cache. Error: %s", err)
		}
		memberID = memberIdentities[0].MemberId
	case consts.IdentityTypePassword:
		// 支持邮箱/手机号作为账号登录
		var memberIdentities []authmodel.MemberIdentity
		for _, v := range []int64{consts.IdentityTypePassword, consts.IdentityTypeEmail, consts.IdentityTypePhone} {
			memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, v, req.Identifier)
			if err != nil {
				return nil, xerror.InternalServerErr
			}
			if len(memberIdentities) > 0 {
				break
			}
		}
		if len(memberIdentities) == 0 {
			l.Errorf("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
			return nil, xerror.NotFoundErr
		}
		if !xstrings.PasswordMatch(memberIdentities[0].Credential, req.Credential) {
			l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
			return nil, xerror.UnauthorizedErr
		}
		memberID = memberIdentities[0].MemberId
	case consts.IdentityTypeWechat:
		// TODO: 实现微信登录
	case consts.IdentityTypeGoogle:
		userInfo, err := google.GetUserInfo(l.ctx, req.Credential)
		if err != nil {
			l.Errorf("Failed to get google user info. Error: %s", err)
			return nil, xerror.InternalServerErr
		}
		thirdPartyId = userInfo.Sub
		memberInfo.Nickname = userInfo.Name
		memberInfo.Avatar = userInfo.Picture
		memberInfo.Email = userInfo.Email
	case consts.IdentityTypeFacebook:
		userInfo, err := facebook.NewFaceBookApp(l.svcCtx.Config.Facebook.AppID, l.svcCtx.Config.Facebook.AppSecret).GetUserInfo(l.ctx, req.Credential)
		if err != nil {
			l.Errorf("Failed to get facebook user info. Error: %s", err)
			return nil, xerror.InternalServerErr
		}
		thirdPartyId = userInfo.UserID
		memberInfo.Email = userInfo.Email
	case consts.IdentityTypeGithub:
		// TODO: 实现github登录
	}

	if len(thirdPartyId) > 0 {
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, req.IdentityType, thirdPartyId)
		if err != nil {
			l.Errorf("Failed to find member identity. Error: %s", err)
			return nil, xerror.InternalServerErr
		}
		if len(memberIdentities) == 0 {
			l.Infof("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
			// 创建新用户
			memberID = uuid.New().String()
			err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
				err = l.svcCtx.MemberIdentityModel.Insert(l.ctx, tx, &authmodel.MemberIdentity{
					MemberId:     memberID,
					IdentityType: req.IdentityType,
					Identifier:   thirdPartyId,
				})
				if err != nil {
					return err
				}
				memberInfo.MemberId = memberID
				memberInfo.Status = consts.MemberStatusEnabled
				err = l.svcCtx.MemberInfoModel.Insert(l.ctx, tx, memberInfo)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				l.Errorf("Failed to create member identity. Error: %s", err)
				return nil, xerror.InternalServerErr
			}
		} else {
			memberID = memberIdentities[0].MemberId
		}
	}

	if memberID == "" {
		l.Errorf("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
		return nil, xerror.NotFoundErr
	}

	// 生成token、refresh token
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
	token, exp, err := jwks.GenerateToken(jwksRecord.Kid, memberID, string(privateKey), l.svcCtx.Config.Auth.AccessExpire, nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}
	refreshToken, refreshExp, err := jwks.GenerateToken(jwksRecord.Kid, memberID, string(privateKey), l.svcCtx.Config.Auth.RefreshExpire, nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}
	refreshTokenHash, err := xstrings.AesCbcEncryptBase64(token, "refresh_token_se", nil)
	if err != nil {
		l.Errorf("Failed to encrypt refresh token. Error: %s", err)
		return nil, err
	}

	// 保存token、refresh token
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = l.svcCtx.RefreshTokenModel.Insert(l.ctx, tx, &authmodel.RefreshToken{
			MemberId:         memberID,
			RefreshTokenHash: refreshTokenHash,
			Status:           consts.RefreshTokenStatusEnabled,
			ExpiredAt:        exp,
		})
		if err != nil {
			return err
		}
		l.svcCtx.TokenRefreshRecordModel.Insert(l.ctx, tx, &authmodel.TokenRefreshRecord{
			RefreshTokenHash: refreshTokenHash,
			Token:            token,
			Kid:              jwksRecord.Kid,
			Ip:               xcontext.GetClientIp(l.ctx),
			ExpiredAt:        refreshExp,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Errorf("Failed to insert token refresh record. Error: %s", err)
		return nil, xerror.InternalServerErr
	}
	resp = &types.MemberLoginResp{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return
}
