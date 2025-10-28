package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"
)

type ContextMiddleware struct {
}

func NewContextMiddleware() *ContextMiddleware {
	return &ContextMiddleware{}
}

func (m *ContextMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if language := r.Header.Get(consts.ACCEPT_LANGUAGE); len(language) > 0 {
			r = r.WithContext(xcontext.WithLang(metadata.AppendToOutgoingContext(r.Context(), consts.ACCEPT_LANGUAGE, language), language))
		}

		addr := httpx.GetRemoteAddr(r)
		// if header has xff, address format: 183.220.108.39,34.128.159.221
		if strings.Contains(addr, ",") {
			addr = strings.Split(addr, ",")[0]
		}
		ip, _, _ := net.SplitHostPort(addr)
		if ip == "" {
			ip = addr
		}
		r = r.WithContext(xcontext.WithClientIP(metadata.AppendToOutgoingContext(r.Context(), consts.CONTEXT_CLIENT_IP, ip), ip))

		next(w, r)
	}
}
