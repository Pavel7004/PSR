all: lint

build:
	@echo "------------------"
	@echo "Building app...   "
	@echo "------------------"
	go build cmd/app/app.go

lint:
	@echo "------------------"
	@echo "Running linter... "
	@echo "------------------"
	golangci-lint run ./...

test:
	@echo "------------------"
	@echo "Running tests... "
	@echo "------------------"
	go test ./... -coverprofile=cover.out

jaeger:
	docker run -dp 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

clear:
	rm app *.out

clearCache:
	go clean -testcache
	go clean -cache

.PHONY: all build swag clean jaeger lint test
