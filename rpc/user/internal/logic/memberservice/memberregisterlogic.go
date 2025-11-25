package memberservicelogic

import (
	"context"
	"fmt"
	"time"

	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/model"
	"github.com/geekeryy/api-hub/rpc/user/user"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MemberRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRegisterLogic {
	return &MemberRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MemberRegisterLogic) MemberRegister(in *user.MemberRegisterReq) (*user.MemberRegisterResp, error) {
	var err error
	memberUUID := uuid.New().String()

	switch in.IdentityType {
	case consts.IdentityTypeEmail:
		err = l.registerEmail(in, memberUUID)
	case consts.IdentityTypePhone:
		err = l.registerPhone(in, memberUUID)
	case consts.IdentityTypePassword:
		err = l.registerPassword(in, memberUUID)
	}
	if err != nil {
		return nil, err
	}

	// 生成token、refresh token
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
	token, exp, err := jwks.GenerateToken(jwksRecord.Kid, memberUUID, privateKey, int64(l.svcCtx.Config.Jwt.AccessExpire), nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}
	refreshToken, refreshExp, err := jwks.GenerateToken(jwksRecord.Kid, memberUUID, privateKey, int64(l.svcCtx.Config.Jwt.RefreshExpire), nil)
	if err != nil {
		l.Errorf("Failed to generate token. Error: %s", err)
		return nil, err
	}
	refreshTokenHash, err := xstrings.AesCbcEncryptBase64(token, l.svcCtx.Config.Secret.RefreshToken, nil)
	if err != nil {
		l.Errorf("Failed to encrypt refresh token. Error: %s", err)
		return nil, err
	}

	// 保存token、refresh token
	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err = l.svcCtx.RefreshTokenModel.Insert(l.ctx, session, &model.RefreshToken{
			MemberId:         memberUUID,
			RefreshTokenHash: refreshTokenHash,
			Status:           consts.RefreshTokenStatusEnabled,
			ExpiredAt:        exp,
		})
		if err != nil {
			return err
		}
		_, err = l.svcCtx.TokenRefreshRecordModel.Insert(l.ctx, session, &model.TokenRefreshRecord{
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
	resp := &user.MemberRegisterResp{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func (l *MemberRegisterLogic) registerPassword(in *user.MemberRegisterReq, memberUUID string) error {
	if len(in.Credential) == 0 {
		return fmt.Errorf("密码不能为空")
	}
	for _, v := range []int64{consts.IdentityTypePassword, consts.IdentityTypeEmail, consts.IdentityTypePhone} {
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, v, in.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) > 0 {
			return fmt.Errorf("账号已注册")
		}
	}
	hash, err := xstrings.PasswordHash(in.Credential)
	if err != nil {
		return err
	}

	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err = l.svcCtx.MemberIdentityModel.Insert(l.ctx, session, &model.MemberIdentity{
			MemberUuid:   memberUUID,
			IdentityType: consts.IdentityTypePassword,
			Identifier:   in.Identifier,
			Credential:   hash,
		})
		if err != nil {
			return err
		}
		_, err = l.svcCtx.MemberInfoModel.Insert(l.ctx, session, &model.MemberInfo{
			MemberUuid: memberUUID,
			Status:     consts.MemberStatusEnabled,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *MemberRegisterLogic) registerPhone(in *user.MemberRegisterReq, memberUUID string) error {
	cacheValue, err := l.svcCtx.RedisClient.Get(l.ctx, fmt.Sprintf("phone_code_%s", in.Identifier)).Result()
	if err != nil {
		return err
	}
	if in.Code != cacheValue {
		return fmt.Errorf("手机验证码不正确")
	}
	memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePhone, in.Identifier)
	if err != nil {
		return err
	}
	if len(memberIdentities) > 0 {
		return fmt.Errorf("手机号已注册")
	}
	insertIdentities := []model.MemberIdentity{{
		MemberUuid:   memberUUID,
		IdentityType: consts.IdentityTypePhone,
		Identifier:   in.Identifier,
	}}
	if len(in.Credential) > 0 {
		hash, err := xstrings.PasswordHash(in.Credential)
		if err != nil {
			return err
		}
		insertIdentities = []model.MemberIdentity{{
			MemberUuid:   memberUUID,
			IdentityType: consts.IdentityTypePhone,
			Identifier:   in.Identifier,
			Credential:   hash,
		}}
		memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePassword, in.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) == 0 {
			insertIdentities = append(insertIdentities, model.MemberIdentity{
				MemberUuid:   memberUUID,
				IdentityType: consts.IdentityTypePassword,
				Identifier:   in.Identifier,
				Credential:   hash,
			})
		}
	}
	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.MemberIdentityModel.BatchInsert(l.ctx, session, insertIdentities)
		if err != nil {
			return err
		}
		_, err = l.svcCtx.MemberInfoModel.Insert(l.ctx, session, &model.MemberInfo{
			MemberUuid: memberUUID,
			Status:     consts.MemberStatusEnabled,
			Phone:      in.Identifier,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *MemberRegisterLogic) registerEmail(in *user.MemberRegisterReq, memberUUID string) error {
	cacheValue, err := l.svcCtx.RedisClient.Get(l.ctx, fmt.Sprintf("email_code_%s", in.Identifier)).Result()
	if err != nil {
		return err
	}
	if in.Code != cacheValue {
		return fmt.Errorf("邮箱验证码不正确")
	}

	memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypeEmail, in.Identifier)
	if err != nil {
		return err
	}
	if len(memberIdentities) > 0 {
		return fmt.Errorf("邮箱已注册")
	}

	insertIdentities := []model.MemberIdentity{{
		MemberUuid:   memberUUID,
		IdentityType: consts.IdentityTypeEmail,
		Identifier:   in.Identifier,
	}}
	if len(in.Credential) > 0 {
		hash, err := xstrings.PasswordHash(in.Credential)
		if err != nil {
			return err
		}
		insertIdentities = []model.MemberIdentity{{
			MemberUuid:   memberUUID,
			IdentityType: consts.IdentityTypeEmail,
			Identifier:   in.Identifier,
			Credential:   hash,
		}}
		memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePassword, in.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) == 0 {
			insertIdentities = append(insertIdentities, model.MemberIdentity{
				MemberUuid:   memberUUID,
				IdentityType: consts.IdentityTypePassword,
				Identifier:   in.Identifier,
				Credential:   hash,
			})
		}
	}
	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.MemberIdentityModel.BatchInsert(l.ctx, session, insertIdentities)
		if err != nil {
			return err
		}
		_, err = l.svcCtx.MemberInfoModel.Insert(l.ctx, session, &model.MemberInfo{
			MemberUuid: memberUUID,
			Status:     consts.MemberStatusEnabled,
			Email:      in.Identifier,
			Birthday:   time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
