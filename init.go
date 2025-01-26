// Package otelgodcp provides OpenTelemetry-based tracing implementations for the go-dcp package.
// This allows users to leverage OpenTelemetry for distributed tracing in their go-dcp applications.
package otelgodcp

import (
	"os"

	"github.com/Trendyol/go-dcp/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// init registers the OpenTelemetry tracer with the go-dcp tracing system.
// This function is automatically invoked when the package is imported, performing the necessary
// context registration steps to enable tracing.
//
// Usage:
//
// To use this package in your project, import it anonymously (with the blank identifier `_`), similar
// to how you import database/sql driver packages. This ensures the init function is executed and
// the OpenTelemetry tracer is registered.
//
// Example:
//
//	```
//	import (
//		_ "github.com/Trendyol/otel-go-dcp"
//	)
//	```
//
// By registering the OpenTelemetry tracer, this package helps integrate OpenTelemetry's powerful
// tracing capabilities with go-dcp, facilitating enhanced observability and monitoring for your
// distributed applications.
func init() {
	// Create a new Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		panic(err)
	}

	// Create a new tracer provider with the exporter and a resource
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(envOr("GO-DCP_COLLECTOR_SERVICE_NAME", "go-dcp")),
		)),
	)

	otel.SetTracerProvider(tp)
	requestTracer := NewOpenTelemetryRequestTracer(tp)
	traceRegisterErr := tracing.RegisterRequestTracer(requestTracer)

	if traceRegisterErr != nil {
		panic(traceRegisterErr)
	}
}

func envOr(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
