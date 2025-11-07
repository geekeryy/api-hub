// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package explorer

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/explorer"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RSS代理
func RssProxyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RssProxyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr)
			return
		}
		if err := svcCtx.Validator.ValidateStruct(r.Context(), req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr.WithMessage(err.Error()))
			return
		}

		l := explorer.NewRssProxyLogic(r.Context(), svcCtx)
		resp, err := l.RssProxy(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
