all: server

server:
	mkdir -p release
	CGO_ENABLED=0 go build -o release/SPaaS_server ./server

javascript:
	mkdir -p release
	CGO_ENABLED=0 go build -o release/javascript ./buildpacks/javascript

server_linux:
	mkdir -p release
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o release/SPaaS_server ./server

frontend:
	npm run build --prefix ./frontend

frontend_deps:
	cd frontend; npm install

release_dev:
	docker build -t mgranderath/spaas:dev .

release:
	docker build -t mgranderath/spaas .

test:
	go test ./... -v

fmt:
	go fmt ./... -v

.PHONY: server frontend release release_dev test fmt
