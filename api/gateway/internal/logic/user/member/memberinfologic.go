package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/geekeryy/api-hub/rpc/user/client/memberservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type MemberInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberInfoLogic {
	return &MemberInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberInfoLogic) MemberInfo() (resp *types.MemberInfoResp, err error) {
	memberUUID := xcontext.GetMemberID(l.ctx)
	if memberUUID == "" {
		return nil, xerror.ForbiddenErr
	}
	memberInfo, err := l.svcCtx.MemberService.GetMemberInfo(l.ctx, &memberservice.GetMemberInfoReq{
		MemberUuid: memberUUID,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.MemberInfoResp{
		Nickname: memberInfo.Nickname,
		Avatar:   memberInfo.Avatar,
		Gender:   int(memberInfo.Gender),
		Birthday: memberInfo.Birthday,
		Phone:    memberInfo.Phone,
		Email:    memberInfo.Email,
	}

	return
}
