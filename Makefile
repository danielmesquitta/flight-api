include .env

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

.PHONY: default
default: run

.PHONY: install
install:
	@go mod download && ./bin/install.sh

.PHONY: update
update:
	@go get -u ./... && go mod tidy

.PHONY: run
run:
	@air -c .air.toml

.PHONY: clear
clear:
	@find ./tmp -mindepth 1 ! -name '.gitkeep' -delete

.PHONY: generate
generate:
	@go generate ./...

.PHONY: docs
docs:
	@swag init -g ./cmd/server/main.go -o ./docs -q && swag2op init -g cmd/server/main.go -q --openapiOutputDir ./tmp && mv ./tmp/swagger.json ./docs/openapi.json && mv ./tmp/swagger.yaml ./docs/openapi.yaml

.PHONY: build
build:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/server ./cmd/server

.PHONY: lint
lint:
	@golangci-lint run && golines **/*.go -m 80 --dry-run

.PHONY: lint-fix
lint-fix:
	@golangci-lint run --fix && golines **/*.go -w -m 80

.PHONY: unit-test
unit-test:
	@ENVIRONMENT=test go test -cover -coverprofile=tmp/coverage.out ./internal/domain/usecase/... ./internal/pkg/... -timeout 5s

.PHONY: integration-test
integration-test:
	@ENVIRONMENT=test go test ./test/integration/... -timeout 3m

.PHONY: test
test: unit-test integration-test
	@true

.PHONY: coverage
coverage:
	@go tool cover -html=tmp/coverage.out

%::
	@true
