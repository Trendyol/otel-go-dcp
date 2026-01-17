package otelgodcp

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func TestServiceAttributesDetector_Detect(t *testing.T) {
	// Mock hostname detector
	mockHostname := "mock-host"
	mockHostnameError := errors.New("hostname error")

	tests := []struct {
		name             string
		opts             Options
		hostnameDetector func() (string, error)
		expectError      bool
		expectedResource *resource.Resource
		expectedAttrs    map[string]string
	}{
		{
			name: "Success with all options provided",
			opts: Options{
				ServiceName:       "my-service",
				ServiceNamespace:  "my-namespace",
				ServiceInstanceID: "instance-1",
				ServiceVersion:    "1.0.0",
			},
			hostnameDetector: func() (string, error) { return mockHostname, nil },
			expectError:      false,
			expectedAttrs: map[string]string{
				string(semconv.ServiceNameKey):       "my-service",
				string(semconv.ServiceNamespaceKey):  "my-namespace",
				string(semconv.ServiceInstanceIDKey): "instance-1",
				string(semconv.ServiceVersionKey):    "1.0.0",
			},
		},
		{
			name: "Success with defaults (missing optional fields)",
			opts: Options{
				ServiceName: "my-service",
			},
			hostnameDetector: func() (string, error) { return mockHostname, nil },
			expectError:      false,
			expectedAttrs: map[string]string{
				string(semconv.ServiceNameKey):       "my-service",
				string(semconv.ServiceNamespaceKey):  NotAvailable,
				string(semconv.ServiceInstanceIDKey): mockHostname,
				string(semconv.ServiceVersionKey):    NotAvailable,
			},
		},
		{
			name: "Success with hostname error (fallback to N/A)",
			opts: Options{
				ServiceName: "my-service",
			},
			hostnameDetector: func() (string, error) { return "", mockHostnameError },
			expectError:      false,
			expectedAttrs: map[string]string{
				string(semconv.ServiceNameKey):       "my-service",
				string(semconv.ServiceNamespaceKey):  NotAvailable,
				string(semconv.ServiceInstanceIDKey): NotAvailable,
				string(semconv.ServiceVersionKey):    NotAvailable,
			},
		},
		{
			name: "Failure missing service name",
			opts: Options{
				ServiceName: "",
			},
			hostnameDetector: func() (string, error) { return mockHostname, nil },
			expectError:      true,
			expectedAttrs:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewServiceAttributesDetector(tt.opts, tt.hostnameDetector)
			res, err := detector.Detect(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				for key, val := range tt.expectedAttrs {
					attr := res.Set().ToSlice()
					found := false
					for _, a := range attr {
						if string(a.Key) == key {
							assert.Equal(t, val, a.Value.AsString(), "Attribute %s should match", key)
							found = true
							break
						}
					}
					assert.True(t, found, "Attribute %s should be present", key)
				}
			}
		})
	}
}
