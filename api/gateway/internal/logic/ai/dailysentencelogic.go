package ai

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailySentenceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 每日一句
func NewDailySentenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailySentenceLogic {
	return &DailySentenceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailySentenceLogic) DailySentence(req *types.DailySentenceReq) (resp *types.DailySentenceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
