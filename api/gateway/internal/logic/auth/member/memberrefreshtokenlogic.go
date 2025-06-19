package member

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberRefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新Token
func NewMemberRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRefreshTokenLogic {
	return &MemberRefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberRefreshTokenLogic) MemberRefreshToken(req *types.MemberRefreshTokenReq) (resp *types.MemberRefreshTokenResp, err error) {
	// todo: add your logic here and delete this line

	return
}
