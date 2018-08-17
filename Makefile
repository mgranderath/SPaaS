all: test build

build:
	mkdir -p release
	go build -o release/SPaaS_server ./server

test:
	go test ./... -v

fmt:
	go fmt ./... -v

.PHONY: build test fmt