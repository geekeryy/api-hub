package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/geekeryy/api-hub/api/oms/internal/logic/auth"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
)

// 登出
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
