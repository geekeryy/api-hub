package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/geekeryy/api-hub/core/language"
	"github.com/geekeryy/api-hub/library/localization"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type baseResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// ErrorHandler 处理http错误消息
// 用法：
//
//	直接返回内置错误，隐藏下层错误
//	   return xerror.SystemErr
//	返回grpc错误，使用下层的内置错误（从details中查找最近一个内置错误）
//	   return grpcerr
func ErrorHandler(ctx context.Context, err error) (statusCode int, errResponse any) {
	if err == nil {
		return http.StatusOK, nil
	}
	target := &xerror.Error{}
	if errors.As(err, &target) {
		logx.Errorw(target.MessageId, logx.Field("callers", target.Slacks), logx.Field("code", target.Code))
		return transform(ctx, target)
	}

	if s, ok := status.FromError(err); ok {
		found := false
		fields := make([]logx.LogField, 0)
		for i := len(s.Details()) - 1; i >= 0; i-- {
			if errors.As(s.Details()[i].(error), &target) {
				fields = append(fields, logx.Field(fmt.Sprintf("%d", i), fmt.Sprintf("%v err:%v", getCaller(target.Slacks), target.MessageId)))
				if !found {
					statusCode, errResponse = transform(ctx, target)
					found = true
				}
			}
		}
		logx.Errorw(s.Message(), fields...)
		if found {
			return
		}
	} else {
		logx.Error(err.Error())
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		return transform(ctx, xerror.InvalidParameterErr.WithStatus(http.StatusBadRequest))
	}

	return http.StatusInternalServerError, baseResponse{
		Code: 500,
		Msg:  "unknown error",
	}
}

func getCaller(slacks []string) string {
	if len(slacks) == 0 {
		return ""
	}
	return slacks[0]
}

// transform 将内置错误转换为http返回
func transform(ctx context.Context, target *xerror.Error) (int, any) {
	statusCode := http.StatusOK
	if target.Status > 0 {
		statusCode = int(target.Status)
	} else if target.Code > 0 && target.Code <= 600 {
		statusCode = int(target.Code)
	}
	messageId := localization.Localize(language.Lang(ctx), &i18n.LocalizeConfig{
		MessageID:    target.MessageId,
		TemplateData: target.TemplateDate,
		PluralCount:  target.Plural,
	})
	return statusCode, baseResponse{
		Code: target.Code,
		Msg:  messageId,
	}
}

// grpcCodeToHTTPStatus 将 gRPC 错误码转换为 HTTP 状态码
func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func OkHandler(ctx context.Context, v interface{}) any {
	return baseResponse{
		Code: 0,
		Msg:  "success",
		Data: v,
	}
}
