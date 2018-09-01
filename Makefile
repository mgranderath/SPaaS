all: build

build:
	mkdir -p release
	go build -o release/SPaaS_server ./server

build_linux:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/SPaaS_server ./server

test:
	go test ./... -v

fmt:
	go fmt ./... -v

.PHONY: build test fmt