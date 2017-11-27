all:
	mkdir -p build
	GOOS=linux GOARCH=arm GOARM=6 go build -o build/Piaas_arm .
	go build -o build/Piaas .