.PHONY: all format test build

all: format test build

format:
	go fmt ./...

test:
	go test -v ./...

build:
	go build -o powerdns-exporter.exe .
