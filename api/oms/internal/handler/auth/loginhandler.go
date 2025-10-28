package auth

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/oms/internal/logic/auth"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"
	"github.com/geekeryy/api-hub/api/oms/internal/types"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 登录
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr)
			return
		}
		if err := svcCtx.Validator.ValidateStruct(r.Context(), req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr.WithMessage(err.Error()))
			return
		}

		l := auth.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
