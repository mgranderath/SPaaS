SRC = $(find -name "*.go" -not -path "./vendor/*")

all: build

.PHONY: build
build: 
	mkdir -p build
	GOOS=linux GOARCH=arm GOARM=6 go build -o build/PiaaS_ARM ./server
	go build -o build/PiaaS ./server

build_rest:
	mkdir -p build
	go build -o build/PiaaS ./server

dependencies:
	glide install --strip-vendor

clean: 
	rm -rf build