module github.com/Trendyol/otel-go-dcp

go 1.20

replace github.com/Trendyol/go-dcp => ../go-dcp

replace go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.20.0

replace go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.20.0

replace go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.20.0

replace go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.20.0

retract (
	v1.2.17
	v1.2.16
)

require (
	github.com/Trendyol/go-dcp v1.2.0-rc.4
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/sdk v1.20.0
	go.opentelemetry.io/otel/trace v1.24.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
