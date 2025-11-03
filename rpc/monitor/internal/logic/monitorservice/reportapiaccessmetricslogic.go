package monitorservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/rpc/monitor/internal/svc"
	"github.com/geekeryy/api-hub/rpc/monitor/monitor"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportApiAccessMetricsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportApiAccessMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportApiAccessMetricsLogic {
	return &ReportApiAccessMetricsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportApiAccessMetricsLogic) ReportApiAccessMetrics(in *monitor.ReportApiAccessMetricsReq) (*monitor.Empty, error) {
	// todo: add your logic here and delete this line

	return &monitor.Empty{}, nil
}
