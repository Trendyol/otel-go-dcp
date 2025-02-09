package otelgodcp

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
)

type Options struct {
	OTLPHeaders              map[string]string `env:"OTEL_EXPORTER_OTLP_HEADERS,overwrite,separator==" json:"OTLPHeaders"`
	ServiceName              string            `env:"OTEL_SERVICE_NAME,overwrite,default=otel-go-dcp" json:"serviceName"`
	ServiceNamespace         string            `env:"OTEL_SERVICE_NAMESPACE,overwrite,default=otel-go-dcp" json:"serviceNamespace"`
	ServiceInstanceID        string            `env:"OTEL_SERVICE_INSTANCE_ID,overwrite,default=otel-go-dcp" json:"serviceInstanceID"`
	ServiceVersion           string            `env:"OTEL_SERVICE_VERSION,overwrite,default=N/A" json:"serviceVersion"`
	OTLPEndpoint             string            `env:"OTEL_EXPORTER_OTLP_ENDPOINT,overwrite,default=http://localhost:4317" json:"OTLPEndpoint"`
	OTLPCompression          string            `env:"OTEL_EXPORTER_OTLP_COMPRESSION,overwrite,default=gzip" json:"OTLPCompression" `
	TraceSamplingProbability float64           `env:"OTEL_TRACES_SAMPLER_ARG,overwrite,default=0.1" json:"TraceSamplingProbability"`
	OTLPTimeout              time.Duration     `env:"OTEL_EXPORTER_OTLP_TIMEOUT,overwrite,default=10s" json:"OTLPTimeout"`
}

func (opts *Options) toResources() (*resource.Resource, error) {
	return resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithDetectors(NewServiceAttributesDetector(*opts, os.Hostname)),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithOS(),
		resource.WithTelemetrySDK(),
	)
}
