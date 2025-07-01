package middleware

import (
	"net/http"
	"strings"

	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/golang-jwt/jwt/v5/keyfunc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type OmsJwtMiddleware struct {
	kfunc keyfunc.Keyfunc
}

func NewOmsJwtMiddleware(kfunc keyfunc.Keyfunc) *OmsJwtMiddleware {
	return &OmsJwtMiddleware{
		kfunc: kfunc,
	}
}

func (m *OmsJwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logx.WithContext(r.Context())
		tokenStr := r.Header.Get(consts.TOKEN)
		if len(tokenStr) > 6 && strings.ToUpper(tokenStr[0:7]) == "BEARER " {
			tokenStr = tokenStr[7:]
		}

		if len(tokenStr) == 0 {
			httpx.ErrorCtx(r.Context(), w, xerror.UnauthorizedErr)
			return
		}

		claims, errV := jwks.ValidateToken(tokenStr, m.kfunc)
		if errV != nil {
			logger.Errorf("validate token :%s error: %v", tokenStr, errV)
			httpx.ErrorCtx(r.Context(), w, xerror.UnauthorizedErr)
			return
		}

		memberId, err := claims.GetSubject()
		if err != nil || len(memberId) == 0 {
			logx.Errorf("validate claims err:%v %+v", err, claims)
			httpx.ErrorCtx(r.Context(), w, xerror.UnauthorizedErr)
			return
		}

		r = r.WithContext(xcontext.WithMemberID(r.Context(), memberId))
		next(w, r)
	}
}
