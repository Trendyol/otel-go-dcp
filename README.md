# otel-go-dcp

`otel-go-dcp` provides OpenTelemetry-based tracing implementations for the `go-dcp` package. This allows users to leverage OpenTelemetry for distributed tracing in their `go-dcp` applications.

## Features

- Integrates OpenTelemetry with `go-dcp` for enhanced observability.
- Automatically registers the OpenTelemetry tracer with the `go-dcp` tracing system.
- Supports Jaeger exporter for tracing data.

## Installation

To install the package, use the following command:

```sh
go get github.com/Trendyol/otel-go-dcp
```

## Usage

To use this package in your project, import it anonymously (with the blank identifier `_`), similar to how you import `database/sql` driver packages. This ensures the `init` function is executed and the OpenTelemetry tracer is registered.

Example:

```go
import (
    _ "github.com/Trendyol/otel-go-dcp"
)
```

By registering the OpenTelemetry tracer, this package helps integrate OpenTelemetry's powerful tracing capabilities with `go-dcp`, facilitating enhanced observability and monitoring for your distributed applications.

## Environment Variables

The following environment variables can be set to configure the tracing behavior:

- `GO-DCP_COLLECTOR_SERVICE_NAME`: The service name to be used by the Jaeger exporter. Defaults to `go-dcp` if not set.
- `OTEL_EXPORTER_JAEGER_ENDPOINT`: The Jaeger collector endpoint. Defaults to `http://localhost:14268/api/traces` if not set.
- `OTEL_EXPORTER_JAEGER_USER`: The username to be used for authentication with the Jaeger collector.
- `OTEL_EXPORTER_JAEGER_PASSWORD`: The password to be used for authentication with the Jaeger collector.

Example:

```sh
export GO-DCP_COLLECTOR_SERVICE_NAME=my-service
export OTEL_EXPORTER_JAEGER_ENDPOINT=http://localhost:14268/api/traces
export OTEL_EXPORTER_JAEGER_USER=my-username
export OTEL_EXPORTER_JAEGER_PASSWORD=my-password
```

Example:

```sh
export GO-DCP_COLLECTOR_SERVICE_NAME=my-service
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
