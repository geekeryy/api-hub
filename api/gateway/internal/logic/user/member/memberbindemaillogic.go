package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberBindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定邮箱
func NewMemberBindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberBindEmailLogic {
	return &MemberBindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberBindEmailLogic) MemberBindEmail(req *types.MemberBindEmailReq) error {
	// todo: add your logic here and delete this line

	return nil
}
