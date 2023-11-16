package crypto

import (
	"encoding/binary"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20"
)

type Oracle struct {
	LocationKey []byte
	Cipher      *AEAD
}

func NewOracle(key []byte) (*Oracle, error) {
	locationKey, err := deriveLocationKey(key)
	if err != nil {
		return nil, err
	}
	cipherKey, err := deriveEncryptionKey(key)
	if err != nil {
		return nil, err
	}
	cipher, err := NewCipher(cipherKey)
	if err != nil {
		return nil, err
	}
	return &Oracle{
		LocationKey: locationKey,
		Cipher:      cipher,
	}, nil
}

func (o *Oracle) LocationKeyFromCounter(counter int64) ([]byte, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, counter)
	return o.LocationKeyFromBytes(buf[:n])
}

func (o *Oracle) LocationKeyFromBytes(id []byte) ([]byte, error) {
	h, err := blake2b.New512(o.LocationKey)
	if err != nil {
		return nil, err
	}
	if _, err := h.Write(id); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func deriveLocationKey(key []byte) ([]byte, error) {
	return chacha20.HChaCha20(key, []byte{0xbf, 0xc8, 0x79, 0xdf, 0x8f, 0x4a, 0xdb, 0xac, 0x6a, 0x43, 0x18, 0xe0, 0x09, 0x26, 0x3d, 0x0d})
}

func deriveEncryptionKey(key []byte) ([]byte, error) {
	return chacha20.HChaCha20(key, []byte{0xa3, 0x3f, 0xac, 0x20, 0x5d, 0x05, 0x7a, 0xa8, 0x3b, 0x71, 0xf2, 0x10, 0x48, 0x2e, 0xb4, 0x26})
}
