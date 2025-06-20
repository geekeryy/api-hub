package admin

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/user/admin"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户信息
func AdminInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminInfoLogic(r.Context(), svcCtx)
		resp, err := l.AdminInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
