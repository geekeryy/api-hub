// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package jwks

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/geekeryy/api-hub/api/oms/internal/logic/jwks"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
)

// jwks密钥列表
func GetJWKSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := jwks.NewGetJWKSLogic(r.Context(), svcCtx)
		resp, err := l.GetJWKS()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
