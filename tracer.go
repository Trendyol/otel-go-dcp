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

const InstrumentationName = "com.trendyol/otel-go-dcp"

type OpenTelemetryRequestTracer struct {
	wrapped trace.Tracer
}

func NewOpenTelemetryRequestTracer(provider trace.TracerProvider) *OpenTelemetryRequestTracer {
	return &OpenTelemetryRequestTracer{
		wrapped: provider.Tracer(InstrumentationName),
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

func (span *OpenTelemetryRequestSpan) SetAttribute(key string, value interface{}) {
	if val := attributeFromValue(key, value); val.Valid() {
		span.wrapped.SetAttributes(val)
		return
	}

	// Handle special cases/fallbacks
	switch v := value.(type) {
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

//nolint:gocyclo
func attributeFromValue(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case *string:
		return attribute.String(key, *v)
	case bool:
		return attribute.Bool(key, v)
	case *bool:
		return attribute.Bool(key, *v)
	case int:
		return attribute.Int(key, v)
	case *int:
		return attribute.Int(key, *v)
	case int64:
		return attribute.Int64(key, v)
	case *int64:
		return attribute.Int64(key, *v)
	case uint32:
		return attribute.Int(key, int(v))
	case *uint32:
		return attribute.Int(key, int(*v))
	case float64:
		return attribute.Float64(key, v)
	case *float64:
		return attribute.Float64(key, *v)
	case []string:
		return attribute.StringSlice(key, v)
	case []bool:
		return attribute.BoolSlice(key, v)
	case []int:
		return attribute.IntSlice(key, v)
	case []int64:
		return attribute.Int64Slice(key, v)
	case []float64:
		return attribute.Float64Slice(key, v)
	}
	return attribute.KeyValue{}
}

func (span *OpenTelemetryRequestSpan) AddEvent(key string, timestamp time.Time) {
	span.wrapped.AddEvent(key, trace.WithTimestamp(timestamp))
}
