package middleware

import (
	"net/http"
	"strings"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JwtMiddleware struct {
	logx.Logger
	kfunc keyfunc.Keyfunc
}

func NewJwtMiddleware(kfunc keyfunc.Keyfunc) *JwtMiddleware {
	return &JwtMiddleware{
		kfunc: kfunc,
	}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logx.WithContext(r.Context())
		tokenStr := stripBearerPrefixFromTokenString(r.Header.Get(consts.TOKEN))
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

		ctx := xcontext.WithKID(r.Context(), claims[jwkset.HeaderKID].(string))
		ctx = xcontext.WithMemberUUID(ctx, memberId)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func stripBearerPrefixFromTokenString(tok string) string {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:]
	}
	return tok
}
