package middleware

import (
	"net/http"

	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/pquerna/otp/totp"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type OmsOtpMiddleware struct {
	secret string
}

func NewOmsOtpMiddleware(secret string) *OmsOtpMiddleware {
	return &OmsOtpMiddleware{
		secret: secret,
	}
}

func (m *OmsOtpMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(consts.OTP_TOKEN)
		if token == "" {
			httpx.Error(w, xerror.UnauthorizedErr)
			return
		}
		if !totp.Validate(token, m.secret) {
			httpx.Error(w, xerror.ForbiddenErr)
			return
		}
		next(w, r)
	}
}
