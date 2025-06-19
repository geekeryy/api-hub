package member

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/auth/member"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 登出
func MemberLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberLogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := svcCtx.Validator.ValidateStruct(r.Context(), req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := member.NewMemberLogoutLogic(r.Context(), svcCtx)
		err := l.MemberLogout(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
