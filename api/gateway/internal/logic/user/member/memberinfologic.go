package member

import (
	"context"
	"time"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/xcontext"

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
	memberInfo, err := l.svcCtx.MemberInfoModel.FindOneByMemberId(l.ctx, xcontext.GetMemberID(l.ctx))
	if err != nil {
		l.Errorf("Failed to find member info. Error: %s", err)
		return nil, err
	}
	resp = &types.MemberInfoResp{
		Nickname: memberInfo.Nickname,
		Avatar:   memberInfo.Avatar,
		Gender:   int(memberInfo.Gender),
		Birthday: memberInfo.Birthday.Format(time.DateTime),
		Phone:    memberInfo.Phone,
		Email:    memberInfo.Email,
	}

	return
}
