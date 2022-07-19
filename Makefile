all: lint

build:
	@echo "------------------"
	@echo "Building app...   "
	@echo "------------------"
	go build cmd/server/server.go

lint:
	@echo "------------------"
	@echo "Running linter... "
	@echo "------------------"
	golangci-lint run ./...

jaeger:
	docker run -dp 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

clear:
	rm shop *.out

clean:
	go clean -testcache
	go clean -cache

.PHONY: all build swag clean jaeger lint
