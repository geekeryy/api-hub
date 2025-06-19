package jwks

import (
	"context"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除公钥
func NewDeleteKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKeyLogic {
	return &DeleteKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteKeyLogic) DeleteKey(req *types.DeleteKeyReq) error {
	ok, err := l.svcCtx.Jwkset.KeyDelete(l.ctx, req.Kid)
	if err != nil {
		l.Errorf("Failed to delete the key. Error: %s", err)
		return err
	}
	if !ok {
		l.Errorf("Key not found. Error: %s", err)
		return err
	}
	return nil
}
