package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	// 第三方平台 会话标识
	ThirdSessionKey = attribute.Key("third.session")

	// 支付 支付意向ID
	PaymentIntendIDKey = attribute.Key("payment.id")

	// 汇率
	RateMapKey = attribute.Key("rate.map")
)

var tracer = otel.GetTracerProvider().Tracer("core/tracing")

func Tracing(ctx context.Context, spanName string, f func(ctx context.Context) error, attrs ...attribute.KeyValue) error {
	var err error
	newCtx := Start(ctx, spanName, attrs...)
	defer End(newCtx, err)
	err = f(newCtx)
	if err != nil {
		return err
	}
	return nil
}

func Start(ctx context.Context, spanName string, attrs ...attribute.KeyValue) context.Context {
	newCtx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindInternal))
	span.SetAttributes(attrs...)
	return newCtx
}

func End(newCtx context.Context, err error, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(newCtx)
	if !span.IsRecording() {
		return
	}
	defer span.End()
	span.SetAttributes(attrs...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

type Attribute struct {
	Key   attribute.Key
	Value string
}
