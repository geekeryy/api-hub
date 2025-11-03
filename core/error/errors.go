package error

import (
	"context"
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc  --go_out=../ error.proto

//
// 错误处理原则：
//
// 1. 隐藏底层错误时必须立即打印详细错误日志，防止日志丢失导致问题难以排查
// 2. 需要展现给用户的错误信息必须使用内置错误，可自定义错误信息，禁止直接返回原始错误
// 3. 错误发生处必须携带错误堆栈信息，便于排查问题
//
// 错误处理方法：
//
// 1. 自身产生的错误使用xerror.New()封装错误信息，返回内置错误
// 2. 来自下游的错误，简单调用可以依赖下游返回的内置错误，复杂组合调用建议处理错误并返回内置错误
//
// 用法：
//	http层：直接返回内置错误，隐藏下层错误
//	grpc层：返回内置错误或普通错误（如果返回内置错误，ErrorInterceptor会将内置错误转换为标准grpc错误）
//
//			return xerror.SystemErr                                    // 直接返回内置错误消息（错误细节需要在发生处打印）
//			return xerror.InvalidParameterErr.WithMessage(err.Error()) // 返回内置错误自定义错误消息
//			return xerror.New(err,xerror.InvalidParameterErr)          // （推荐）将内置错误追加到原始错误之后，给上层兜底（如果上层直接抛出错误，则http.errorHandler将负责打印详细日志以及抛出内置错误）
//			return err                                                 // 当前层不做处理，直接返回原始错误

func (e *Error) Error() string {
	return e.MessageId
}

func (e *Error) clone() *Error {
	return &Error{
		Code:         e.Code,
		MessageId:    e.MessageId,
		TemplateDate: e.TemplateDate,
		Plural:       e.Plural,
		Slacks:       e.Slacks,
		Status:       e.Status,
		Details:      e.Details,
	}
}

func (e *Error) WithSlacks() *Error {
	if len(e.Slacks) == 0 {
		return e.clone().setSlacks()
	}
	return e
}

func (e *Error) WithDetails(details ...string) *Error {
	return e.WithSlacks().setDetails(details...)
}

func (e *Error) WithMessage(message string) *Error {
	return e.WithSlacks().setMessage(message)
}

func (e *Error) WithMetadata(key, value string) *Error {
	return e.WithSlacks().setMetadata(key, value)
}

func (e *Error) WithPlural() *Error {
	return e.WithSlacks().setPlural(2)
}

func (e *Error) setSlacks() *Error {
	e.Slacks = callers()
	return e
}

func (e *Error) setDetails(details ...string) *Error {
	e.Details = append(e.Details, details...)
	return e
}

func (e *Error) setMessage(message string) *Error {
	e.MessageId = message
	return e
}
func (e *Error) setMetadata(key, value string) *Error {
	e.TemplateDate[key] = value
	return e
}
func (e *Error) setPlural(plural int64) *Error {
	e.Plural = plural
	return e
}

func Callers() []string {
	return callers()
}

// callers 获取堆栈信息
func callers() []string {
	pc := make([]uintptr, 3)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[1:n])
	slacks := make([]string, 0)
	for {
		f, more := frames.Next()
		slacks = append(slacks, fmt.Sprintf("%s:%d", f.File, f.Line))
		if !more {
			break
		}
	}
	return slacks
}

// ErrorInterceptor 将内置错误转换为标准grpc错误
func ErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		var coreErr *Error
		if errors.As(err, &coreErr) {
			s, _ := status.New(codes.Unknown, coreErr.MessageId).WithDetails(coreErr)
			return resp, s.Err()
		}
	}
	return resp, err
}
