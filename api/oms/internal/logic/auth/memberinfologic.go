package auth

import (
	"context"

	"github.com/geekeryy/api-hub/api/oms/internal/svc"
	"github.com/geekeryy/api-hub/api/oms/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
