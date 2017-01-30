.DEFAULT_GOAL := all

deps:
	go get golang.org/x/crypto/ssh/terminal
	go get golang.org/x/crypto/nacl/secretbox
	go get golang.org/x/crypto/scrypt
	go get github.com/boltdb/bolt

test:
	go test -v ./...

build:
	go build -v .

all: deps test build
