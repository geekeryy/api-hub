package memberservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/model/usermodel"
	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberInfoLogic {
	return &GetMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMemberInfoLogic) GetMemberInfo(in *user.GetMemberInfoReq) (*user.GetMemberInfoResp, error) {
	memberInfo, err := l.svcCtx.MemberInfoModel.FindOneByMemberId(l.ctx, in.MemberId)
	if err != nil {
		l.Errorf("Failed to find member info. Error: %s, member_id: %s", err, in.MemberId)
		if errors.Is(err, usermodel.ErrNotFound) {
			return nil, xerror.NotFoundErr.WithMetadata("member_id", in.MemberId)
		}
		return nil, xerror.DBErr.WithSlacks()
	}
	resp := &user.GetMemberInfoResp{
		MemberId: memberInfo.MemberId,
		Nickname: memberInfo.Nickname,
		Avatar:   memberInfo.Avatar,
		Gender:   memberInfo.Gender,
		Birthday: memberInfo.Birthday.Format(time.DateTime),
		Phone:    memberInfo.Phone,
		Email:    memberInfo.Email,
	}

	return resp, nil
}
