package jwks

import (
	"context"
	"time"

	"github.com/geekeryy/api-hub/api/oms/internal/svc"
	"github.com/geekeryy/api-hub/api/oms/internal/types"

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
	records, err := l.svcCtx.TokenRefreshRecordModel.FindByKid(l.ctx, req.Kid)
	if err != nil {
		l.Errorf("Failed to find token refresh record by kid. Error: %s", err)
		return nil, err
	}

	resp = &types.KeyUsageResp{
		Records: make([]types.TokenRefreshRecord, 0),
	}
	groupByTime := make(map[string]int)
	for _, record := range records {
		timeStr := record.CreatedAt.Format(time.DateTime)
		groupByTime[timeStr]++
	}
	for timeStr, count := range groupByTime {
		resp.Records = append(resp.Records, types.TokenRefreshRecord{
			Time:  timeStr,
			Count: count,
		})
	}
	return
}
