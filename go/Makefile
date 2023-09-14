all: tools build run

tools:
	go generate -tags=tools ./...

build:
	go generate ./...
	go work sync > /dev/null 2>&1 || true
	go build ./...

run:
	go run cmd/kmttg.go

.PHONY: all tools tidy build run
