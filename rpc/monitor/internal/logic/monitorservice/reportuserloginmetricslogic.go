package monitorservicelogic

import (
	"context"

	"github.com/geekeryy/api-hub/rpc/monitor/internal/svc"
	"github.com/geekeryy/api-hub/rpc/monitor/monitor"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportUserLoginMetricsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportUserLoginMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportUserLoginMetricsLogic {
	return &ReportUserLoginMetricsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportUserLoginMetricsLogic) ReportUserLoginMetrics(in *monitor.ReportUserLoginMetricsReq) (*monitor.Empty, error) {
	for _, item := range in.Items {
		l.Infof("ReportUserLoginMetricsItem: %+v", item)
	}
	return &monitor.Empty{}, nil
}
