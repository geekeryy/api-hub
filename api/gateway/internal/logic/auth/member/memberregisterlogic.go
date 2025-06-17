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

type MemberRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册
func NewMemberRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRegisterLogic {
	return &MemberRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TODO: 高并发情况下，会导致重复注册，需要优化
func (l *MemberRegisterLogic) MemberRegister(req *types.MemberRegisterReq) error {
	switch req.IdentityType {
	case consts.IdentityTypeEmail:
		return l.registerEmail(req)
	case consts.IdentityTypePhone:
		return l.registerPhone(req)
	case consts.IdentityTypePassword:
		return l.registerPassword(req)
	}
	return nil
}

func (l *MemberRegisterLogic) registerPassword(req *types.MemberRegisterReq) error {
	if len(req.Password) == 0 {
		return fmt.Errorf("密码不能为空")
	}
	for _, v := range []int64{consts.IdentityTypePassword, consts.IdentityTypeEmail, consts.IdentityTypePhone} {
		memberIdentities, err := l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, v, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) > 0 {
			return fmt.Errorf("账号已注册")
		}
	}
	memberID := uuid.New().String()
	hash, err := xstrings.PasswordHash(req.Password)
	if err != nil {
		return err
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = l.svcCtx.MemberIdentityModel.Insert(l.ctx, tx, &authmodel.MemberIdentity{
			MemberId:     memberID,
			IdentityType: consts.IdentityTypePassword,
			Identifier:   req.Identifier,
			Credential:   hash,
		})
		if err != nil {
			return err
		}
		err = l.svcCtx.MemberInfoModel.Insert(l.ctx, tx, &usermodel.MemberInfo{
			MemberId: memberID,
			Status:   consts.MemberStatusEnabled,
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

func (l *MemberRegisterLogic) registerPhone(req *types.MemberRegisterReq) error {
	cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("phone_code_%s", req.Identifier))
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
	if len(memberIdentities) > 0 {
		return fmt.Errorf("手机号已注册")
	}
	memberID := uuid.New().String()
	insertIdentities := []authmodel.MemberIdentity{{
		MemberId:     memberID,
		IdentityType: consts.IdentityTypePhone,
		Identifier:   req.Identifier,
	}}
	if len(req.Password) > 0 {
		hash, err := xstrings.PasswordHash(req.Password)
		if err != nil {
			return err
		}
		insertIdentities = []authmodel.MemberIdentity{{
			MemberId:     memberID,
			IdentityType: consts.IdentityTypePhone,
			Identifier:   req.Identifier,
			Credential:   hash,
		}}
		memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePassword, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) == 0 {
			insertIdentities = append(insertIdentities, authmodel.MemberIdentity{
				MemberId:     memberID,
				IdentityType: consts.IdentityTypePassword,
				Identifier:   req.Identifier,
				Credential:   hash,
			})
		}
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = l.svcCtx.MemberIdentityModel.BatchInsert(l.ctx, tx, insertIdentities)
		if err != nil {
			return err
		}
		err = l.svcCtx.MemberInfoModel.Insert(l.ctx, tx, &usermodel.MemberInfo{
			MemberId: memberID,
			Status:   consts.MemberStatusEnabled,
			Phone:    req.Identifier,
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

func (l *MemberRegisterLogic) registerEmail(req *types.MemberRegisterReq) error {
	cacheValue, err := l.svcCtx.Cache.Get(fmt.Sprintf("email_code_%s", req.Identifier))
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
	if len(memberIdentities) > 0 {
		return fmt.Errorf("邮箱已注册")
	}

	memberID := uuid.New().String()
	insertIdentities := []authmodel.MemberIdentity{{
		MemberId:     memberID,
		IdentityType: consts.IdentityTypeEmail,
		Identifier:   req.Identifier,
	}}
	if len(req.Password) > 0 {
		hash, err := xstrings.PasswordHash(req.Password)
		if err != nil {
			return err
		}
		insertIdentities = []authmodel.MemberIdentity{{
			MemberId:     memberID,
			IdentityType: consts.IdentityTypeEmail,
			Identifier:   req.Identifier,
			Credential:   hash,
		}}
		memberIdentities, err = l.svcCtx.MemberIdentityModel.FindByIdentity(l.ctx, consts.IdentityTypePassword, req.Identifier)
		if err != nil {
			return err
		}
		if len(memberIdentities) == 0 {
			insertIdentities = append(insertIdentities, authmodel.MemberIdentity{
				MemberId:     memberID,
				IdentityType: consts.IdentityTypePassword,
				Identifier:   req.Identifier,
				Credential:   hash,
			})
		}
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = l.svcCtx.MemberIdentityModel.BatchInsert(l.ctx, tx, insertIdentities)
		if err != nil {
			return err
		}
		err = l.svcCtx.MemberInfoModel.Insert(l.ctx, tx, &usermodel.MemberInfo{
			MemberId: memberID,
			Status:   consts.MemberStatusEnabled,
			Email:    req.Identifier,
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
