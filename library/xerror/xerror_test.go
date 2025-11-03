package xerror_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	coreError "github.com/geekeryy/api-hub/core/error"
	"github.com/geekeryy/api-hub/core/handler"
	"github.com/geekeryy/api-hub/core/language"
	"github.com/geekeryy/api-hub/core/xcontext"
	_ "github.com/geekeryy/api-hub/library/localization" // 初始化翻译模块
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestError(t *testing.T) {
	e1 := xerror.NotFoundErr.WithSlacks()
	e2 := xerror.NotFoundErr
	fmt.Println(e1.Slacks)
	fmt.Println(e2.Slacks)

}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			"野生err",
			errors.New("normal error"),
			handler.BaseResponse{
				Code: 500,
				Msg:  "unknown error",
			},
		},
		{
			"标准grpc错误",
			status.Error(codes.Unauthenticated, "unauthenticated"),
			handler.BaseResponse{
				Code: 500,
				Msg:  "unknown grpc error",
			},
		},
		{
			"业务错误",
			xerror.NotFoundErr.WithSlacks(),
			handler.BaseResponse{
				Code: 404,
				Msg:  "资源不存在",
			},
		},
		{
			"多条错误信息",
			xerror.NotFoundErr.WithDetails("detail1", "detail2"),
			handler.BaseResponse{
				Code: 404,
				Msg:  "资源不存在",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := xcontext.WithLang(context.Background(), language.ZH)
			if _, got := handler.ErrorHandler(logx.WithContext(context.Background()))(ctx, tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorHandler(err) got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorInterceptor(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			"业务错误",
			xerror.NotFoundErr.WithSlacks(),
			handler.BaseResponse{
				Code: 404,
				Msg:  "资源不存在",
			},
		},
		{
			"多条错误信息",
			xerror.NotFoundErr.WithDetails("detail1", "detail2"),
			handler.BaseResponse{
				Code: 404,
				Msg:  "资源不存在",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := xcontext.WithLang(context.Background(), language.ZH)
			_, grpcgot := coreError.ErrorInterceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, tt.err
			})
			if _, got := handler.ErrorHandler(logx.WithContext(context.Background()))(ctx, grpcgot); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorHandler(err) got %v, want %v", got, tt.want)
			}
		})
	}
}
