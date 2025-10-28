package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/geekeryy/api-hub/api/oms/internal/logic/auth"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
)

// 获取用户信息
func MemberInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewMemberInfoLogic(r.Context(), svcCtx)
		resp, err := l.MemberInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
