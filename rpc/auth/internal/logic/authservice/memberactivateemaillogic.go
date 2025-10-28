package authservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/library/consts"
	"github.com/geekeryy/api-hub/rpc/auth/auth"
	"github.com/geekeryy/api-hub/rpc/auth/internal/svc"
	"github.com/geekeryy/api-hub/rpc/model/membermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberActivateEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberActivateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberActivateEmailLogic {
	return &MemberActivateEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MemberActivateEmailLogic) MemberActivateEmail(in *auth.MemberActivateEmailReq) (*auth.Empty, error) {
	identities, err := l.svcCtx.MemberIdentityModel.FindByMemberUUID(l.ctx, in.MemberUuid)
	if err != nil {
		return nil, err
	}

	identityPassword := ""
	for _, identity := range identities {
		if identity.IdentityType == consts.IdentityTypePassword {
			identityPassword = identity.Credential
		}
	}
	_, err = l.svcCtx.MemberIdentityModel.Insert(l.ctx, nil, &membermodel.MemberIdentity{
		MemberUuid:   in.MemberUuid,
		IdentityType: consts.IdentityTypeEmail,
		Identifier:   in.Email,
		Credential:   identityPassword,
	})
	if err != nil {
		return nil, err
	}

	return &auth.Empty{}, nil
}
