package healthz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/healthz"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
)

// 健康检查
func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := healthz.NewPingLogic(r.Context(), svcCtx)
		err := l.Ping()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
