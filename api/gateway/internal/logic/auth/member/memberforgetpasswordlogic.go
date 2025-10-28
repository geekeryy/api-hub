package member

import (
	"context"
	"fmt"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/rpc/model/membermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberForgetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 忘记密码
func NewMemberForgetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberForgetPasswordLogic {
	return &MemberForgetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberForgetPasswordLogic) MemberForgetPassword(req *types.MemberForgetPasswordReq) error {
	var memberIdentitie *membermodel.MemberIdentity
	switch req.IdentityType {
	case consts.IdentityTypeEmail:
		cacheValue, err := l.svcCtx.RedisClient.Get(l.ctx, fmt.Sprintf("email_code_%s", req.Identifier)).Result()
		if err != nil {
			return err
		}
		if req.Code != cacheValue {
			return fmt.Errorf("邮箱验证码不正确")
		}
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypeEmail, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) > 1 {
			return fmt.Errorf("邮箱绑定多个账号")
		}
		if len(memberIdentities) == 0 {
			return fmt.Errorf("邮箱未绑定")
		}
		memberIdentitie = memberIdentities[0]
	case consts.IdentityTypePhone:
		cacheValue, err := l.svcCtx.RedisClient.Get(l.ctx, fmt.Sprintf("phone_code_%s", req.Identifier)).Result()
		if err != nil {
			return err
		}
		if req.Code != cacheValue {
			return fmt.Errorf("手机验证码不正确")
		}
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePhone, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) > 1 {
			return fmt.Errorf("手机号绑定多个账号")
		}
		if len(memberIdentities) == 0 {
			return fmt.Errorf("手机号未绑定")
		}
		memberIdentitie = memberIdentities[0]
	default:
		return fmt.Errorf("身份类型不正确")
	}

	hash, err := xstrings.PasswordHash(req.Password)
	if err != nil {
		return err
	}

	err = l.svcCtx.MemberIdentityModel.UpdateCredential(l.ctx, memberIdentitie.Id, hash)
	if err != nil {
		return err
	}
	return nil
}
