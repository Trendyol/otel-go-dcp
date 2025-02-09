package otelgodcp

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Trendyol/go-dcp/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type OpenTelemetryRequestTracer struct {
	wrapped trace.Tracer
}

func NewOpenTelemetryRequestTracer(provider trace.TracerProvider) *OpenTelemetryRequestTracer {
	return &OpenTelemetryRequestTracer{
		wrapped: provider.Tracer("com.trendyol/otel-go-dcp"),
	}
}

func (tracer *OpenTelemetryRequestTracer) RequestSpan(parentContext tracing.RequestSpanContext, operationName string) tracing.RequestSpan {
	ctx, span := tracer.wrapped.Start(parentContext.RefCtx, operationName)
	parentContext.RefCtx = ctx
	return NewOpenTelemetryRequestSpan(parentContext, span)
}

type OpenTelemetryRequestSpan struct {
	ctx     tracing.RequestSpanContext
	wrapped trace.Span
}

func NewOpenTelemetryRequestSpan(ctx tracing.RequestSpanContext, span trace.Span) *OpenTelemetryRequestSpan {
	return &OpenTelemetryRequestSpan{
		ctx:     ctx,
		wrapped: span,
	}
}

func (span *OpenTelemetryRequestSpan) End() {
	span.wrapped.End()
}

func (span *OpenTelemetryRequestSpan) Context() tracing.RequestSpanContext {
	return span.ctx
}

//nolint:funlen
func (span *OpenTelemetryRequestSpan) SetAttribute(key string, value interface{}) {
	switch v := value.(type) {
	case string:
		span.wrapped.SetAttributes(attribute.String(key, v))
	case *string:
		span.wrapped.SetAttributes(attribute.String(key, *v))
	case bool:
		span.wrapped.SetAttributes(attribute.Bool(key, v))
	case *bool:
		span.wrapped.SetAttributes(attribute.Bool(key, *v))
	case int:
		span.wrapped.SetAttributes(attribute.Int(key, v))
	case *int:
		span.wrapped.SetAttributes(attribute.Int(key, *v))
	case int64:
		span.wrapped.SetAttributes(attribute.Int64(key, v))
	case *int64:
		span.wrapped.SetAttributes(attribute.Int64(key, *v))
	case uint32:
		span.wrapped.SetAttributes(attribute.Int(key, int(v)))
	case *uint32:
		span.wrapped.SetAttributes(attribute.Int(key, int(*v)))
	case float64:
		span.wrapped.SetAttributes(attribute.Float64(key, v))
	case *float64:
		span.wrapped.SetAttributes(attribute.Float64(key, *v))
	case []string:
		span.wrapped.SetAttributes(attribute.StringSlice(key, v))
	case []bool:
		span.wrapped.SetAttributes(attribute.BoolSlice(key, v))
	case []int:
		span.wrapped.SetAttributes(attribute.IntSlice(key, v))
	case []int64:
		span.wrapped.SetAttributes(attribute.Int64Slice(key, v))
	case []float64:
		span.wrapped.SetAttributes(attribute.Float64Slice(key, v))
	case fmt.Stringer:
		span.wrapped.SetAttributes(attribute.String(key, v.String()))
	case map[string]interface{}, struct{}, interface{}:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			log.Printf("Unable to marshal object for key %s: %v", key, err)
		} else {
			span.wrapped.SetAttributes(attribute.String(key, string(jsonBytes)))
		}
	default:
		log.Println("Unable to determine value as a type that we can handle")
	}
}

func (span *OpenTelemetryRequestSpan) AddEvent(key string, timestamp time.Time) {
	span.wrapped.AddEvent(key, trace.WithTimestamp(timestamp))
}
