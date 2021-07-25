.PHONY: test, lint, run

build:
	@echo todo

install:
	@go install -mod=vendor ./...

test:
	@go test -mod=vendor -race ./...

lint:
	@echo todo

run:
	go run cmd/nuker/main.go