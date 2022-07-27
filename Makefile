.PHONY: build
build:
	@echo todo

.PHONY: install
install:
	@echo "🖥️  Installing Nuker..."
	@go install -mod=vendor ./...
	@echo "✅ Success"

.PHONY: test
test:
	@go test -mod=vendor -race ./...

.PHONY: lint
lint:
	@golangci-lint run ./... --config ./build/golangci-lint/config.yaml

.PHONY: run
run:
	@go run cmd/nuker/main.go

.PHONY: clear
clear:
	@rm -f nuker-*.jsonl

.PHONY: fmt
fmt:
	@gofumpt -l -w -extra .
