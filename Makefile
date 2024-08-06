all: build

build: swag
	go build -o _output/web-service-gin main.go

run: swag
	go run main.go

swag:
	swag init

deps:
	go get ./...

clean:
	rm -rf _output

.PHONY: all clean build run swag deps
