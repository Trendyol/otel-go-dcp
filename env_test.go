package otelgodcp

import (
	"testing"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

func TestToResources_AlignsWithSDKSchema(t *testing.T) {
	opts := Options{
		ServiceName:      "saas-content-indexer",
		ServiceNamespace: "discovery",
		ServiceInstanceID: "instance-1",
		ServiceVersion:   "1.0.0",
	}

	res, err := opts.toResources()
	require.NoError(t, err)
	require.NotNil(t, res)

	attrs := res.Attributes()
	require.NotEmpty(t, attrs)

	schemaURL := res.SchemaURL()
	require.True(t,
		schemaURL == "" || schemaURL == semconv.SchemaURL,
		"unexpected schema URL: %s", schemaURL,
	)
}

func TestServiceAttributesDetector_UsesSemconvV139(t *testing.T) {
	opts := Options{
		ServiceName:       "test-service",
		ServiceNamespace:  "test-ns",
		ServiceInstanceID: "test-id",
		ServiceVersion:    "v1",
	}

	detector := NewServiceAttributesDetector(opts, func() (string, error) {
		return "hostname", nil
	})

	res, err := detector.Detect(t.Context())
	require.NoError(t, err)
	require.Equal(t, semconv.SchemaURL, res.SchemaURL())
}
