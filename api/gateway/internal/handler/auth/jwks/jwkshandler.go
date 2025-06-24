package jwks

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/auth/jwks"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
)

// 获取公钥
func JWKSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := jwks.NewJWKSLogic(r.Context(), svcCtx)
		err := l.JWKS(w)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
