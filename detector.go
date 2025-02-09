package otelgodcp

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const NotAvailable = "N/A"

type serviceAttributesDetector struct {
	hostnameDetector func() string
	opts             Options
}

// Detect returns a *resource.Resource that describes the service version info.
func (s *serviceAttributesDetector) Detect(_ context.Context) (*resource.Resource, error) {
	serviceNamespace := strings.TrimSpace(s.opts.ServiceNamespace)
	if serviceNamespace == "" {
		serviceNamespace = NotAvailable
	}

	serviceName := strings.TrimSpace(s.opts.ServiceName)
	if serviceName == "" {
		return nil, fmt.Errorf("failed to detect service information, couldn't find service name")
	}

	serviceInstanceID := strings.TrimSpace(s.opts.ServiceInstanceID)
	if serviceInstanceID == "" {
		serviceInstanceID = s.hostnameDetector()
	}
	if serviceInstanceID == "" {
		serviceInstanceID = NotAvailable
	}

	serviceVersion := strings.TrimSpace(s.opts.ServiceVersion)
	if serviceVersion == "" {
		serviceVersion = NotAvailable
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceNamespaceKey.String(serviceNamespace),
		semconv.ServiceInstanceIDKey.String(serviceInstanceID),
		semconv.ServiceVersionKey.String(serviceVersion),
	), nil
}

func NewServiceAttributesDetector(opts Options, hostnameDetector func() (string, error)) resource.Detector {
	return &serviceAttributesDetector{
		opts: opts,
		hostnameDetector: func() string {
			hostname, err := hostnameDetector()
			if err != nil {
				return ""
			}
			return hostname
		},
	}
}
