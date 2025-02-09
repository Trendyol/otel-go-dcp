package otelgodcp

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace"
	otel "go.opentelemetry.io/otel/trace"
)

var _ trace.Sampler = DeterministicSampler{}

type DeterministicSampler struct {
	innerSampler trace.Sampler
}

func NewDeterministicSampler(fraction float64) DeterministicSampler {
	var innerSampler trace.Sampler
	switch {
	case fraction <= 0:
		innerSampler = trace.NeverSample()
	case fraction >= 1:
		innerSampler = trace.AlwaysSample()
	default:
		innerSampler = trace.ParentBased(trace.TraceIDRatioBased(fraction))
	}
	return DeterministicSampler{
		innerSampler: innerSampler,
	}
}

func (d DeterministicSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	result := d.innerSampler.ShouldSample(parameters)
	if result.Decision != trace.RecordAndSample {
		return trace.SamplingResult{
			Decision:   trace.RecordOnly,
			Tracestate: otel.SpanContextFromContext(parameters.ParentContext).TraceState(),
		}
	}
	return result
}

func (d DeterministicSampler) Description() string {
	return fmt.Sprintf("DeterministicSampler{innerSampler:%s}",
		d.innerSampler.Description(),
	)
}
