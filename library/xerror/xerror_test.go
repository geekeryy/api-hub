package xerror_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/geekeryy/api-hub/core/handler"
	"github.com/geekeryy/api-hub/core/language"
	"github.com/geekeryy/api-hub/core/xcontext"
	_ "github.com/geekeryy/api-hub/library/localization" // 初始化翻译模块
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestError(t *testing.T) {
	err := errors.New("normal error")
	err = errors.Wrap(err, "attempt to wrap")
	xerror.New(err)
	fmt.Println(err)
}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			"normal error",
			errors.New("normal error"),
			handler.BaseResponse{
				Code: 500,
				Msg:  "unknown error",
			},
		},
		{
			"standard grpc error",
			status.Error(codes.Unauthenticated, "unauthenticated"),
			handler.BaseResponse{
				Code: 500,
				Msg:  "unknown error",
			},
		},
		{
			"grpc error",
			xerror.InternalServerErr.Rpc(),
			handler.BaseResponse{
				Code: 500,
				Msg:  "系统错误",
			},
		},
		{
			"warp normal error to grpc error",
			xerror.New(errors.New("normal error")),
			handler.BaseResponse{
				Code: int64(codes.Unknown),
				Msg:  "normal error",
			},
		},
		{
			"append details",
			xerror.New(errors.New("normal error"), xerror.InternalServerErr),
			handler.BaseResponse{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := xcontext.WithLang(context.Background(), language.ZH)
			if _, got := handler.ErrorHandler(ctx, tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorHandler(err) got %v, want %v", got, tt.want)
			}
		})
	}
}
