package otelgodcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestNewDeterministicSampler(t *testing.T) {
	tests := []struct {
		name        string
		fraction    float64
		description string
	}{
		{
			name:        "Always Sample (>= 1)",
			fraction:    1.0,
			description: "AlwaysOnSampler",
		},
		{
			name:        "Never Sample (<= 0)",
			fraction:    0.0,
			description: "AlwaysOffSampler",
		},
		{
			name:        "Ratio Sample (0.5)",
			fraction:    0.5,
			description: "TraceIDRatioBased{0.5}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := NewDeterministicSampler(tt.fraction)
			assert.NotNil(t, sampler)
			// We can check the description to verify the inner sampler type partly
			desc := sampler.Description()
			assert.Contains(t, desc, "DeterministicSampler")
			assert.Contains(t, desc, tt.description)
		})
	}
}

func TestDeterministicSampler_ShouldSample(t *testing.T) {
	// Test that it delegates correctly
	sampler := NewDeterministicSampler(1.0) // Always sample

	params := trace.SamplingParameters{
		ParentContext: context.Background(),
		TraceID:       oteltrace.TraceID{1},
		Name:          "test-span",
		Kind:          oteltrace.SpanKindServer,
	}

	result := sampler.ShouldSample(params)
	assert.Equal(t, trace.RecordAndSample, result.Decision)

	// Test Never Sample
	sampler = NewDeterministicSampler(0.0)
	result = sampler.ShouldSample(params)
	// Expect RecordOnly because of our custom logic:
	// The code says: if result.Decision != trace.RecordAndSample { return trace.SamplingResult{ Decision: trace.RecordOnly ... } }
	// Wait, trace.NeverSample returns Drop. So we expect RecordOnly.
	assert.Equal(t, trace.RecordOnly, result.Decision)
}
