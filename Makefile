SRC = $(find -name "*.go" -not -path "./vendor/*")

all: build

build: check
	mkdir -p build
	GOOS=linux GOARCH=arm GOARM=6 go build -o build/Piaas_arm .
	go build -o build/Piaas .

dependencies:
	glide install --strip-vendor

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done

clean: 
	rm -rf build