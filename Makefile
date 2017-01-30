.DEFAULT_GOAL := all

test:
	go test -v ./...

build:
	go build -v .

all: test build
