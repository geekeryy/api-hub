package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/core/xgrpc"
	"github.com/geekeryy/api-hub/rpc/monitor/client/monitorservice"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type ContextMiddleware struct {
	logx.Logger
	MonitorService monitorservice.MonitorService
	reportChan     chan *monitorservice.ReportApiAccessMetricsItem
}

func NewContextMiddleware(monitorLazyClient *xgrpc.LazyClient, logger logx.Logger) *ContextMiddleware {
	m := &ContextMiddleware{
		Logger:         logger,
		MonitorService: monitorservice.NewMonitorService(monitorLazyClient),
		reportChan:     make(chan *monitorservice.ReportApiAccessMetricsItem, 10000),
	}
	// 启动多个协程并发上报监控数据，提高上报效率
	for i := 0; i < 1; i++ {
		go m.reportApiAccessMetricsLoop()
	}
	return m
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

		startTime := time.Now()

		next(w, r)

		duration := time.Since(startTime).Milliseconds()
		m.report(r, w, ip, duration)
	}
}

func (m *ContextMiddleware) report(r *http.Request, w http.ResponseWriter, ip string, duration int64) {
	code, _ := strconv.Atoi(w.Header().Get("status"))
	kid := xcontext.GetKID(r.Context())
	m.reportChan <- &monitorservice.ReportApiAccessMetricsItem{
		Ip:        ip,
		Kid:       kid,
		Api:       r.URL.Path,
		Status:    int64(code),
		Duration:  duration,
		Timestamp: time.Now().UnixNano(),
		TraceId:   trace.SpanFromContext(r.Context()).SpanContext().TraceID().String(),
	}
}

func (m *ContextMiddleware) reportApiAccessMetricsLoop() {
	metricsLen := 1000 // 缓存1000条监控数据，也是单次上报的数据条数
	metrics, err := lru.New[string, *monitorservice.ReportApiAccessMetricsItem](metricsLen)
	if err != nil {
		m.Logger.WithFields(logx.LogField{Key: "error", Value: err}).Error("Failed to create metrics cache")
		return
	}
	reportMetrics := func() {
		_, err := m.MonitorService.ReportApiAccessMetrics(context.Background(), &monitorservice.ReportApiAccessMetricsReq{Items: metrics.Values()})
		if err != nil {
			m.Logger.WithFields(logx.LogField{Key: "error", Value: err}).Error("Failed to report api access metrics")
		}
	}
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case item := <-m.reportChan:
			if metrics.Len() >= metricsLen {
				reportMetrics()
				metrics.Purge()
			}
			if ok, evicted := metrics.ContainsOrAdd(item.TraceId, item); ok || evicted {
				m.Logger.WithFields(logx.LogField{Key: "trace_id", Value: item.TraceId}).Error("trace_id already exists")
			}
		case <-ticker.C:
			if metrics.Len() > 0 {
				reportMetrics()
				metrics.Purge()
			}
		}
	}
}
