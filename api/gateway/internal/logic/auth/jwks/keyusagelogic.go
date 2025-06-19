package jwks

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type KeyUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 公钥使用记录
func NewKeyUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeyUsageLogic {
	return &KeyUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KeyUsageLogic) KeyUsage(req *types.KeyUsageReq) (resp *types.KeyUsageResp, err error) {
	// todo: add your logic here and delete this line

	return
}
