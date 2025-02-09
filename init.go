// Package otelgodcp provides OpenTelemetry-based tracing implementations for the go-dcp package.
// This allows users to leverage OpenTelemetry for distributed tracing in their go-dcp applications.
package otelgodcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"

	"github.com/Trendyol/go-dcp/tracing"
	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	ctx := context.Background()

	// Initialize options
	opts := Options{}
	if err := envconfig.Process(ctx, &opts); err != nil {
		panic(err)
	}

	// Create a new Jaeger exporter
	exp, err := newTraceExporter(opts)
	if err != nil {
		panic(err)
	}

	bsp := trace.NewBatchSpanProcessor(exp)
	sampler := NewDeterministicSampler(opts.TraceSamplingProbability)

	res, toResErr := opts.toResources()
	if toResErr != nil {
		panic(toResErr)
	}

	// Create a new tracer provider with the exporter and a resource
	tp := trace.NewTracerProvider(
		trace.WithSampler(sampler),
		trace.WithSpanProcessor(bsp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	requestTracer := NewOpenTelemetryRequestTracer(tp)
	traceRegisterErr := tracing.RegisterRequestTracer(requestTracer)

	if traceRegisterErr != nil {
		panic(traceRegisterErr)
	}
}

func newTraceExporter(opts Options) (*otlptrace.Exporter, error) {
	grpcOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithHeaders(opts.OTLPHeaders),
		otlptracegrpc.WithCompressor(opts.OTLPCompression),
		otlptracegrpc.WithDialOption(grpc.WithUserAgent(fmt.Sprintf("%s/grpc-go/%s", opts.ServiceName, grpc.Version))),
	}

	parsedEndpoint, err := url.Parse(opts.OTLPEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to setup OTel exporter, couldn't parse otlp endpoint %s: %w", opts.OTLPEndpoint, err)
	}

	if parsedEndpoint.Scheme == "https" {
		certPool, certPoolErr := x509.SystemCertPool()
		if certPoolErr != nil {
			return nil, fmt.Errorf("failed to setup OTel exporter, couldn't get system cert pool: %w", certPoolErr)
		}
		tlsConfig := &tls.Config{InsecureSkipVerify: false, RootCAs: certPool} //nolint:gosec
		grpcOpts = append(grpcOpts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
	}
	grpcOpts = append(grpcOpts, otlptracegrpc.WithEndpoint(parsedEndpoint.Host))

	client := otlptracegrpc.NewClient(grpcOpts...)
	return otlptrace.New(context.Background(), client)
}
