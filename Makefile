APP_NAME := go-api
OPENAPI_SPEC := api/openapi.yaml
OPENAPI_OUT := internal/gen/openapi/api.gen.go
OPENAPI_PACKAGE := openapi

SQLC_CONFIG := sqlc.yaml

.PHONY: help
help:
	@echo "available targets:"
	@echo "  tool-install       install pinned codegen tools"
	@echo "  generate           run all code generation"
	@echo "  generate-openapi   generate OpenAPI server/types code"
	@echo "  generate-sqlc      generate sqlc code"
	@echo "  tidy               tidy go modules"
	@echo "  fmt                format go code"
	@echo "  vet                run go vet"
	@echo "  test               run tests"
	@echo "  build              build server binary"
	@echo "  run                run application"
	@echo "  clean              remove build artifacts"

.PHONY: tool-install
tool-install:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.6.0
	go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

.PHONY: generate
generate: generate-openapi generate-sqlc

.PHONY: generate-openapi
generate-openapi:
	mkdir -p $$(dirname $(OPENAPI_OUT))
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		-generate types,std-http \
		-package $(OPENAPI_PACKAGE) \
		-o $(OPENAPI_OUT) \
		$(OPENAPI_SPEC)

.PHONY: generate-sqlc
generate-sqlc:
	go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate -f $(SQLC_CONFIG)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) ./cmd/server

.PHONY: run
run:
	go run ./cmd/server

.PHONY: clean
clean:
	rm -rf bin
