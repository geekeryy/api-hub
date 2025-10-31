package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/core/xgrpc"
	"github.com/geekeryy/api-hub/rpc/monitor/client/monitorservice"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"
)

type ContextMiddleware struct {
	MonitorService monitorservice.MonitorService
}

func NewContextMiddleware(monitorLazyClient *xgrpc.LazyClient) *ContextMiddleware {
	return &ContextMiddleware{
		MonitorService: monitorservice.NewMonitorService(monitorLazyClient),
	}
}

// 添加accept-language和clientip到上下文
func (m *ContextMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if language := r.Header.Get(consts.ACCEPT_LANGUAGE); len(language) > 0 {
			r = r.WithContext(xcontext.WithLang(metadata.AppendToOutgoingContext(r.Context(), consts.ACCEPT_LANGUAGE, language), language))
		}

		addr := httpx.GetRemoteAddr(r)
		// if header has xff, address format: 183.220.108.39,34.128.159.221
		if strings.Contains(addr, ",") {
			addr = strings.Split(addr, ",")[0]
		}
		ip, _, _ := net.SplitHostPort(addr)
		if ip == "" {
			ip = addr
		}
		r = r.WithContext(xcontext.WithClientIP(metadata.AppendToOutgoingContext(r.Context(), consts.CONTEXT_CLIENT_IP, ip), ip))

		_, err := m.MonitorService.ReportUserLoginMetrics(r.Context(), &monitorservice.ReportUserLoginMetricsReq{
			Items: []*monitorservice.ReportUserLoginMetricsItem{
				{
					Ip:       ip,
					Service:  "gateway",
					Kid:      "gateway",
					Api:      r.URL.Path,
					Duration: 0,
					Status:   200,
				},
			},
		})
		if err != nil {
			logx.Errorf("Failed to report user login metrics. Error: %s", err)
		}

		next(w, r)
	}
}
