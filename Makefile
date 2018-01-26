SRC = $(find -name "*.go" -not -path "./vendor/*")

all: build

.PHONY: build
build: 
	mkdir -p build
	GOOS=linux GOARCH=arm GOARM=6 go build -o build/PiaaS_ARM ./server
	go build -o build/PiaaS ./server
	go build -o build/PiaaS_cli ./cli

client:
	mkdir -p build
	go build -o build/PiaaS_cli ./cli

rest:
	mkdir -p build
	go build -o build/PiaaS ./server

arm:
	mkdir -p build
	go build -o build/PiaaS_ARM ./server

dependencies:
	glide install --strip-vendor

clean: 
	rm -rf build