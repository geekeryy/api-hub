package member

import (
	"context"
	"fmt"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/geekeryy/api-hub/rpc/model/usermodel"
	"github.com/google/uuid"
	"gorm.io/gorm"

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

// TODO 手机号绑定两个账号的情况，需要优化
func (l *MemberForgetPasswordLogic) MemberForgetPassword(req *types.MemberForgetPasswordReq) error {
	var memberIdentities []authmodel.MemberIdentity
	switch req.IdentityType {
	case consts.IdentityTypeEmail:
		cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("email_code_%s", req.Identifier))
		if err != nil {
			return err
		}
		if req.Code != cacheValue {
			return fmt.Errorf("邮箱验证码不正确")
		}
		memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypeEmail, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) == 0 {
			return fmt.Errorf("邮箱未绑定")
		}
	case consts.IdentityTypePhone:
		cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("phone_code_%s", req.Identifier))
		if err != nil {
			return err
		}
		if req.Code != cacheValue {
			return fmt.Errorf("手机验证码不正确")
		}
		memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePhone, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) == 0 {
			return fmt.Errorf("手机号未绑定")
		}
	default:
		return fmt.Errorf("身份类型不正确")
	}

	hash, err := xstrings.PasswordHash(req.Password)
	if err != nil {
		return err
	}
	memberIdentities[0].Credential = hash

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = l.svcCtx.MemberIdentityModel.Update(l.ctx, tx, &memberIdentities[0])
		if err != nil {
			return err
		}
		// TODO 密码更新记录，禁止使用最近设置的密码，需要优化
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
