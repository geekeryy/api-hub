package member

import (
	"context"
	"fmt"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/facebook"
	"github.com/geekeryy/api-hub/core/google"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/library/xerror"

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
	// TODO 第三方登录，先调用api获取用户信息，再查询数据库
	memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, req.IdentityType, req.Identifier)
	if err != nil {
		l.Errorf("Failed to find member identity. Error: %s", err)
		return nil, xerror.New(xerror.InternalServerErr)
	}

	if len(memberIdentities) == 0 {
		l.Errorf("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
		return nil, xerror.New(xerror.NotFoundErr)
	}

	var memberID string
	for _, memberIdentity := range memberIdentities {
		switch req.IdentityType {
		case consts.IdentityTypePhone:
			cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("member_login_phone_%s", req.Identifier))
			if err != nil {
				l.Errorf("Failed to get cache. Error: %s", err)
				return nil, xerror.New(xerror.InternalServerErr)
			}
			defer l.svcCtx.Cache.Delete(fmt.Sprintf("member_login_phone_%s", req.Identifier))
			if req.Credential != cacheValue {
				l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
				continue
			}
			memberID = memberIdentity.MemberId
		case consts.IdentityTypeEmail:
			cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("member_login_email_%s", req.Identifier))
			if err != nil {
				l.Errorf("Failed to get cache. Error: %s", err)
				return nil, xerror.New(xerror.InternalServerErr)
			}
			defer l.svcCtx.Cache.Delete(fmt.Sprintf("member_login_email_%s", req.Identifier))
			if req.Credential != cacheValue {
				l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
				continue
			}
			memberID = memberIdentity.MemberId
		case consts.IdentityTypePassword:
			if !xstrings.PasswordMatch(memberIdentity.Credential, req.Credential) {
				l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
				continue
			}
			memberID = memberIdentity.MemberId
		case consts.IdentityTypeWechat:
			// TODO: 实现微信登录
			memberID = memberIdentity.MemberId
		case consts.IdentityTypeGoogle:
			userInfo, err := google.GetUserInfo(l.ctx, req.Credential)
			if err != nil {
				l.Errorf("Failed to get google user info. Error: %s", err)
				return nil, xerror.New(xerror.InternalServerErr)
			}
			if userInfo.Sub != memberIdentity.Identifier {
				l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
				continue
			}
			memberID = memberIdentity.MemberId
		case consts.IdentityTypeFacebook:
			userInfo, err := facebook.NewFaceBookApp(l.svcCtx.Config.Facebook.AppID, l.svcCtx.Config.Facebook.AppSecret).GetUserInfo(l.ctx, req.Credential)
			if err != nil {
				l.Errorf("Failed to get facebook user info. Error: %s", err)
				return nil, xerror.New(xerror.InternalServerErr)
			}
			if userInfo.UserID != memberIdentity.Identifier {
				l.Infof("Member login failed. Identity: %s, Credential: %s code not match", req.Identifier, req.Credential)
				continue
			}
			memberID = memberIdentity.MemberId
		case consts.IdentityTypeGithub:
			// TODO: 实现github登录
			memberID = memberIdentity.MemberId
		}
	}

	if memberID == "" {
		l.Errorf("Member identity not found. IdentityType: %d, Identity: %s", req.IdentityType, req.Identifier)
		return nil, xerror.New(xerror.NotFoundErr)
	}

	jwksRecord, err := l.svcCtx.JwksModel.FindLatest(l.ctx)
	if err != nil {
		l.Errorf("Failed to find latest jwks public. Error: %s", err)
		return nil, xerror.New(xerror.InternalServerErr)
	}
	privateKey, err := xstrings.AesCbcDecryptBase64(jwksRecord.PrivateKey, "private_key_secr", nil)
	if err != nil {
		l.Errorf("Failed to decrypt private key. Error: %s", err)
		return nil, err
	}
	token, err := jwks.GenerateToken(jwksRecord.Kid, "1234567890", string(privateKey), l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}
	resp = &types.MemberLoginResp{
		Token:        token,
		RefreshToken: token,
	}

	return
}
