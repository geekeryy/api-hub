package middleware

import (
	"net/http"

	"github.com/MicahParks/keyfunc/v3"
)

type AdminJwtMiddleware struct {
	kfunc keyfunc.Keyfunc
}

func NewAdminJwtMiddleware(kfunc keyfunc.Keyfunc) *AdminJwtMiddleware {
	return &AdminJwtMiddleware{
		kfunc: kfunc,
	}
}

func (m *AdminJwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
