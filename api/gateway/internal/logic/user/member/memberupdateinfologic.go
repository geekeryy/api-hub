package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberUpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户信息
func NewMemberUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberUpdateInfoLogic {
	return &MemberUpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberUpdateInfoLogic) MemberUpdateInfo(req *types.MemberUpdateInfoReq) error {
	// todo: add your logic here and delete this line

	return nil
}
