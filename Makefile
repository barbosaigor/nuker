.PHONY: test, lint, run, clear

build:
	@echo todo

install:
	@echo "🖥️  Installing Nuker..."
	@go install -mod=vendor ./...
	@echo "✅ Success"

test:
	@go test -mod=vendor -race ./...

lint:
	golangci-lint run ./... --config ./build/golangci-lint/config.yaml

run:
	go run cmd/nuker/main.go

clear:
	rm nuker-*.jsonl