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
