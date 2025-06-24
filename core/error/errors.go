package error

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

//go:generate protoc  --go_out=../ error.proto

func (e *Error) Error() string {
	return e.MessageId
}

func (e *Error) WithMessage(message string) *Error {
	e.MessageId = message
	return e
}

func (e *Error) WithStatus(status int64) *Error {
	e.Status = status
	return e
}

func (e *Error) WithMetadata(key, value string) *Error {
	e.TemplateDate[key] = value
	return e
}

func (e *Error) WithPlural() *Error {
	e.Plural = 2
	return e
}

// 将内置错误转换为grpc error
// 用法：
//
//	grpc错误 TODO 用中间件处理？
//		return xerror.SystemErr.Rpc()
//	http错误
//		return xerror.NotFoundErr
func (e *Error) Rpc() error {
	if len(e.Slacks) == 0 {
		e.Slacks = callers()
	}
	s, _ := status.New(codes.Unknown, e.MessageId).WithDetails(e)
	return s.Err()
}

// New 创建一个grpc错误
// 用法：
//
// 如果原始err未附带堆栈信息，则包裹一层堆栈信息
//
//	将普通错误转换为grpc error，来自本层的错误，没有堆栈信息
//	  New(err)
//	将grpc error附加details，来自下层grpc服务的status错误，已有堆栈信息
//	  New(grpcerr, e1, e2)
//
// err: 普通错误则包裹一层堆栈信息，grpc错误则附加details
func New(err error, errs ...protoadapt.MessageV1) error {
	if err == nil {
		return nil
	}

	var slacks []string
	if stackTracer, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		st := stackTracer.StackTrace()
		for _, v := range st {
			b, _ := v.MarshalText()
			slacks = append(slacks, string(b))
		}
	} else {
		slacks = callers()
	}

	s, ok := status.FromError(err)
	if !ok {
		errs = append([]protoadapt.MessageV1{
			&Error{
				Code:      int64(codes.Unknown),
				MessageId: err.Error(),
				Slacks:    slacks,
			},
		}, errs...)
	}
	s, _ = s.WithDetails(errs...)
	return s.Err()
}

// callers 获取堆栈信息
func callers() []string {
	pc := make([]uintptr, 4)
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
