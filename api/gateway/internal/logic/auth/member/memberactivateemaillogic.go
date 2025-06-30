package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberActivateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 激活邮箱
func NewMemberActivateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberActivateEmailLogic {
	return &MemberActivateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberActivateEmailLogic) MemberActivateEmail(req *types.MemberActivateEmailReq) error {
	claims, err := jwks.ValidateToken(req.Token, l.svcCtx.Kfunc)
	if err != nil {
		return err
	}
	l.Infof("claims: %v", claims)
	memberId, err := claims.GetSubject()
	if err != nil {
		return err
	}
	identities, err := l.svcCtx.MemberIdentityModel.FindByMemberId(l.ctx, memberId)
	if err != nil {
		return err
	}
	email, err := jwks.MapClaimsParseString(claims, "email")
	if err != nil {
		return err
	}
	identityPassword := ""
	for _, identity := range identities {
		if identity.IdentityType == consts.IdentityTypePassword {
			identityPassword = identity.Identifier
		}
	}
	err = l.svcCtx.MemberIdentityModel.Insert(l.ctx, nil, &authmodel.MemberIdentity{
		MemberId:     memberId,
		IdentityType: consts.IdentityTypeEmail,
		Identifier:   email,
		Credential:   identityPassword,
	})
	if err != nil {
		return err
	}
	return nil
}
