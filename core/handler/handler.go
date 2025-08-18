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
// 用法：
//
// 1. 返回内置错误
// 2. 打印错误日志
//
//			http层：直接返回内置错误，隐藏下层错误
//			    return xerror.SystemErr                                    // 直接返回内置错误消息（错误细节需要在发生出打印）
//		    	return xerror.InvalidParameterErr.WithMessage(err.Error()) // 返回自定义错误消息（错误细节需要在发生出打印）
//	            return xerror.New(err,xerror.InvalidParameterErr)          // 将内置错误追加到原始错误之后，给下层兜底
//		    	return err                                                 // 当前层不做处理，直接返回原始错误（普通错误、grpc错误）
//			grpc层：使用下层的内置错误（从details中查找最近一个内置错误）
//			    return xerror.SystemErr.Rpc()                     // 直接返回内置错误消息（错误细节需要在发生出打印）
//		    	return xerror.New(err,xerror.InvalidParameterErr) // 将内置错误追加到原始错误之后，给下层兜底
//	                                                              // 将原始错误带给上层，让上层进行逻辑处理
//		    	return err                                        // 当前层不做处理，让上层处理要返回的错误（错误细节需要在发生出打印）
func ErrorHandler(ctx context.Context, err error) (statusCode int, errResponse any) {
	if err == nil {
		return http.StatusOK, nil
	}

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
		logx.WithContext(ctx).Errorw(target.MessageId, fields...)
		return transform(ctx, target)
	}

	s, ok := status.FromError(err)
	if !ok {
		logx.WithContext(ctx).Error(err.Error())
		return http.StatusInternalServerError, BaseResponse{
			Code: 500,
			Msg:  "unknown error",
		}
	}

	found := false
	for _, detail := range s.Proto().Details {
		detail, err := detail.UnmarshalNew()
		if err != nil {
			logx.WithContext(ctx).Errorw("unmarshal detail", logx.Field("error", err))
			continue
		}
		target, ok = detail.(*coreError.Error)
		if !ok {
			logx.WithContext(ctx).Errorw("unknown detail", logx.Field("error", detail))
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
			logx.WithContext(ctx).Errorw(s.Message(), fields...)
			statusCode, errResponse = transform(ctx, target)
			found = true
		}
	}

	if found {
		return statusCode, errResponse
	}

	logx.WithContext(ctx).WithFields(logx.Field("code", uint32(s.Code()))).Error(s)

	return http.StatusInternalServerError, BaseResponse{
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
