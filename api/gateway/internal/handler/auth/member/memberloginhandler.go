package member

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/auth/member"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 登录
func MemberLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := svcCtx.Validator.ValidateStruct(r.Context(), req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := member.NewMemberLoginLogic(r.Context(), svcCtx)
		resp, err := l.MemberLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
