package jwks

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/geekeryy/api-hub/api/oms/internal/logic/jwks"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
)

// 轮换公钥
func RotateKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := jwks.NewRotateKeyLogic(r.Context(), svcCtx)
		err := l.RotateKey()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
