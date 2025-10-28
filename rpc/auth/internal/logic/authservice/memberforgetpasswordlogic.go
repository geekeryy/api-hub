package authservicelogic

import (
	"context"
	"fmt"

	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/rpc/auth/auth"
	"github.com/geekeryy/api-hub/rpc/auth/internal/svc"
	"github.com/geekeryy/api-hub/rpc/model/membermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberForgetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberForgetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberForgetPasswordLogic {
	return &MemberForgetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MemberForgetPasswordLogic) MemberForgetPassword(in *auth.MemberForgetPasswordReq) (*auth.Empty, error) {
	var memberIdentitie *membermodel.MemberIdentity
	switch in.IdentityType {
	case consts.IdentityTypeEmail:
		cacheValue, err := l.svcCtx.RedisClient.Get(l.ctx, fmt.Sprintf("email_code_%s", in.Identifier)).Result()
		if err != nil {
			return nil, err
		}
		if in.Code != cacheValue {
			return nil, fmt.Errorf("邮箱验证码不正确")
		}
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypeEmail, in.Identifier)
		if err != nil {
			return nil, err
		}
		if len(memberIdentities) > 1 {
			return nil, fmt.Errorf("邮箱绑定多个账号")
		}
		if len(memberIdentities) == 0 {
			return nil, fmt.Errorf("邮箱未绑定")
		}
		memberIdentitie = memberIdentities[0]
	case consts.IdentityTypePhone:
		cacheValue, err := l.svcCtx.RedisClient.Get(l.ctx, fmt.Sprintf("phone_code_%s", in.Identifier)).Result()
		if err != nil {
			return nil, err
		}
		if in.Code != cacheValue {
			return nil, fmt.Errorf("手机验证码不正确")
		}
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePhone, in.Identifier)
		if err != nil {
			return nil, err
		}
		if len(memberIdentities) > 1 {
			return nil, fmt.Errorf("手机号绑定多个账号")
		}
		if len(memberIdentities) == 0 {
			return nil, fmt.Errorf("手机号未绑定")
		}
		memberIdentitie = memberIdentities[0]
	default:
		return nil, fmt.Errorf("身份类型不正确")
	}

	hash, err := xstrings.PasswordHash(in.Password)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.MemberIdentityModel.UpdateCredential(l.ctx, memberIdentitie.Id, hash)
	if err != nil {
		return nil, err
	}

	return &auth.Empty{}, nil
}
