all: tools build run

tools:
	go generate -tags=tools ./...

tidy:
	GONOPROXY=https://github.com/tartale/go go mod tidy

build:
	go generate ./...
	go build ./...

run:
	go run cmd/kmttg.go

.PHONY: all tools tidy build run
