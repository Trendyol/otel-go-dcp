.PHONY: default test

default: init

init:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.15.0

clean:
	rm -rf ./build

linter:
	fieldalignment -fix ./...
	golangci-lint run -c .golangci.yml --timeout=5m -v --fix

lint:
	golangci-lint run -c .golangci.yml --timeout=5m -v

test:
	go test ./... .

tidy:
	go mod tidy
	cd test/integration/basic-otel-tracing && go mod tidy

create-tracing-infra:
	bash scripts/create_tracing_infra.sh

delete-tracing-infra:
	bash scripts/delete_tracing_infra.sh
