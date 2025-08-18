package adminservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminInfoLogic {
	return &GetAdminInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminInfoLogic) GetAdminInfo(in *user.GetAdminInfoReq) (*user.GetAdminInfoResp, error) {
	// todo: add your logic here and delete this line

	return &user.GetAdminInfoResp{}, nil
}
