// Package xerror @Description  TODO
// @Author  	 jiangyang
// @Created  	 2024/7/19 下午5:03
package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/geekeryy/api-hub/library/xerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			"normal error",
			errors.New("normal error"),
			baseResponse{
				Code: 500,
				Msg:  "unknown error",
			},
		},
		{
			"standard grpc error",
			status.Error(codes.Unauthenticated, "unauthenticated"),
			baseResponse{
				Code: 500,
				Msg:  "unknown error",
			},
		},
		{
			"grpc error",
			xerror.InternalServerErr.Rpc(),
			baseResponse{
				Code: 500,
				Msg:  "系统错误",
			},
		},
		{
			"warp normal error to grpc error",
			xerror.New(errors.New("normal error")),
			baseResponse{
				Code: int64(codes.Unknown),
				Msg:  "normal error",
			},
		},
		{
			"append details",
			xerror.New(errors.New("normal error"), xerror.InternalServerErr),
			baseResponse{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := ErrorHandler(context.Background(), tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorHandler(err) got %v, want %v", got, tt.want)
			}
		})
	}
}
