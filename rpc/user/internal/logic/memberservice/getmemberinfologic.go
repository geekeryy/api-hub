package memberservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/model"
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
	memberInfo, err := l.svcCtx.MemberInfoModel.FindOneByMemberUuid(l.ctx, in.MemberUuid)
	if err != nil {
		l.Errorf("Failed to find member info. Error: %s, member_uuid: %s", err, in.MemberUuid)
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerror.NotFoundErr.WithMetadata("member_uuid", in.MemberUuid)
		}
		return nil, xerror.DBErr.WithSlacks()
	}
	resp := &user.GetMemberInfoResp{
		MemberUuid: in.MemberUuid,
		Nickname:   memberInfo.Nickname,
		Avatar:     memberInfo.Avatar,
		Gender:     memberInfo.Gender,
		Birthday:   memberInfo.Birthday.Format(time.DateTime),
		Phone:      memberInfo.Phone,
		Email:      memberInfo.Email,
	}

	return resp, nil
}
