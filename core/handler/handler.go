package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	coreError "github.com/geekeryy/api-hub/core/error"
	"github.com/geekeryy/api-hub/core/language"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// ErrorHandler 处理http错误消息
// 1. 处理内置错误，支持国际化
// 2. 处理标准grpc错误，提取最新内置错误
// 3. 处理未知错误，打印详细错误日志
// 4. 返回http响应
func ErrorHandler(logger logx.Logger) func(ctx context.Context, err error) (statusCode int, errResponse any) {
	return func(ctx context.Context, err error) (statusCode int, errResponse any) {
		if err == nil {
			return http.StatusOK, nil
		}
		logger := logger.WithContext(ctx).WithCallerSkip(4)

		// 处理内置错误
		target := &coreError.Error{}
		if errors.As(err, &target) {
			fields := make([]logx.LogField, 0, len(target.Slacks)+1)
			fields = append(fields, logx.Field("code", target.Code))
			for i, v := range target.Slacks {
				fields = append(fields, logx.Field(fmt.Sprintf("callers[%d]", i), v))
			}
			for i, v := range target.Details {
				fields = append(fields, logx.Field(fmt.Sprintf("detail[%d]", i), v))
			}
			if target.OriginalError != "" {
				fields = append(fields, logx.Field("originalError", target.OriginalError))
			}
			logger.Errorw(target.MessageId, fields...)
			return transform(ctx, target)
		}

		s, ok := status.FromError(err)
		if !ok {
			// 未知错误
			logger.WithFields(logx.Field("error", err)).Error("unknown error")
			return http.StatusInternalServerError, BaseResponse{
				Code: 500,
				Msg:  "unknown error",
			}
		}

		// 处理标准grpc错误
		found := false
		proto := s.Proto()
		for i := len(proto.Details) - 1; i >= 0; i-- {
			detail, err := proto.Details[i].UnmarshalNew()
			if err != nil {
				logger.Errorw("UnmarshalNew detail", logx.Field("error", err))
				continue
			}
			target, ok = detail.(*coreError.Error)
			if !ok {
				// 未知细节，打印
				logger.Errorw("unknown detail", logx.Field("error", detail))
				continue
			}
			if !found {
				fields := make([]logx.LogField, 0)
				for i, v := range target.Slacks {
					fields = append(fields, logx.Field(fmt.Sprintf("callers[%d]", i), v))
				}
				for i, v := range target.Details {
					fields = append(fields, logx.Field(fmt.Sprintf("detail[%d]", i), v))
				}
				logger.Errorw(s.Message(), fields...)
				statusCode, errResponse = transform(ctx, target)
				found = true
			}
		}

		if found {
			return statusCode, errResponse
		}

		// 未知grpc错误，http直接返回了没有附加内置错误的grpc错误
		logger.WithFields(logx.Field("code", uint32(s.Code())), logx.Field("err", s)).Error("unknown grpc error")

		return http.StatusInternalServerError, BaseResponse{
			Code: 500,
			Msg:  "unknown grpc error",
		}
	}
}

// transform 将内置错误转换为http返回
func transform(ctx context.Context, target *coreError.Error) (int, any) {
	statusCode := http.StatusOK
	if target.Status > 0 {
		statusCode = int(target.Status)
	} else if target.Code > int64(codes.OK) && target.Code < int64(codes.Unauthenticated) {
		statusCode = grpcCodeToHTTPStatus(codes.Code(target.Code))
	} else if target.Code > 100 && target.Code < 600 {
		statusCode = int(target.Code)
	}
	messageId := language.Localize(language.Lang(ctx), &i18n.LocalizeConfig{
		MessageID:    target.MessageId,
		TemplateData: target.TemplateDate,
		PluralCount:  target.Plural,
	})
	return statusCode, BaseResponse{
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
	return BaseResponse{
		Code: 0,
		Msg:  "success",
		Data: v,
	}
}
