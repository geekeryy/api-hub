package member

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/auth/member"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 注册
func MemberRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberRegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr)
			return
		}
		if err := svcCtx.Validator.ValidateStruct(r.Context(), req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr.WithMessage(err.Error()))
			return
		}

		l := member.NewMemberRegisterLogic(r.Context(), svcCtx)
		err := l.MemberRegister(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
