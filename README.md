# otel-go-dcp

`otel-go-dcp` provides OpenTelemetry-based tracing implementations for the `go-dcp` package. This allows users to leverage OpenTelemetry for distributed tracing in their `go-dcp` applications.

## Features

- Integrates OpenTelemetry with `go-dcp` for enhanced observability.
- Automatically registers the OpenTelemetry tracer with the `go-dcp` tracing system.

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

Here is the updated environment variables table with the `Type` field added and the `Required` column removed:

## Environment Variables

The following environment variables can be set to configure the tracing behavior:

| Variable Name                     | Description                                      | Type     | Default Value          | Example                          |
|-----------------------------------|--------------------------------------------------|----------|------------------------|----------------------------------|
| `OTEL_EXPORTER_OTLP_HEADERS`      | Headers for OTLP exporter                        | `string` |                        | `key1=value1,key2=value2`        |
| `OTEL_SERVICE_NAME`               | Name of the service                              | `string` | `otel-go-dcp`          | `my-service`                     |
| `OTEL_SERVICE_NAMESPACE`          | Namespace of the service                         | `string` | `otel-go-dcp`          | `my-namespace`                   |
| `OTEL_SERVICE_INSTANCE_ID`        | Instance ID of the service                       | `string` | `otel-go-dcp`          | `instance-123`                   |
| `OTEL_SERVICE_VERSION`            | Version of the service                           | `string` | `N/A`                  | `1.0.0`                          |
| `OTEL_EXPORTER_OTLP_ENDPOINT`     | Endpoint for OTLP exporter                       | `string` | `http://localhost:4317`| `http://collector:4317`          |
| `OTEL_EXPORTER_OTLP_COMPRESSION`  | Compression for OTLP exporter                    | `string` | `gzip`                 | `none`                           |
| `OTEL_TRACES_SAMPLER_ARG`         | Trace sampling probability                       | `float`  | `0.1`                  | `0.5`                            |
| `OTEL_EXPORTER_OTLP_TIMEOUT`      | Timeout for OTLP exporter                        | `duration`| `10s`                  | `5s`                             |

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
