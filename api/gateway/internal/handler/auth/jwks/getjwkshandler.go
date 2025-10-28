package jwks

import (
	"net/http"

	"github.com/geekeryy/api-hub/api/gateway/internal/logic/auth/jwks"
	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取公钥
func GetJWKSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JWKSReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr)
			return
		}
		if err := svcCtx.Validator.ValidateStruct(r.Context(), req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerror.InvalidParameterErr.WithMessage(err.Error()))
			return
		}

		l := jwks.NewGetJWKSLogic(r.Context(), svcCtx)
		err := l.GetJWKS(&req, w)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
