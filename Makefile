all:
	go get golang.org/x/crypto/ssh/terminal
	go get golang.org/x/crypto/nacl/secretbox
	go get golang.org/x/crypto/scrypt

	go install
