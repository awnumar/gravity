package main

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

const Overhead int = chacha20poly1305.Overhead + chacha20poly1305.NonceSizeX

type AEAD struct {
	aead cipher.AEAD
}

func NewCipher(key []byte) (*AEAD, error) {
	aead, err := chacha20poly1305.New(key)
	return &AEAD{aead}, err
}

func (c *AEAD) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return c.aead.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *AEAD) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < c.aead.NonceSize() {
		return nil, errors.New("ciphertext too small")
	}
	return c.aead.Open(nil, ciphertext[:c.aead.NonceSize()], ciphertext[c.aead.NonceSize():], nil)
}
