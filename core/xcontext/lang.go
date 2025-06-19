package xcontext

import (
	"context"

	"github.com/geekeryy/api-hub/core/consts"
	"google.golang.org/grpc/metadata"
)

type acceptlanguage struct{}

func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, acceptlanguage{}, lang)
}

func GetLang(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if value := ctx.Value(acceptlanguage{}); value != nil {
			return value.(string)
		}
		return ""
	}
	lang, ok := md[consts.ACCEPT_LANGUAGE]
	if !ok || len(lang) == 0 || len(lang[0]) == 0 {
		return ""
	}
	return lang[0]
}
