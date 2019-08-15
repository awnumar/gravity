package main

import (
	"bytes"
	"testing"

	"github.com/awnumar/memguard"
)

func TestEncryptDecrypt(t *testing.T) {
	// Declare the plaintext and the key.
	m := make([]byte, 64)
	memguard.ScrambleBytes(m)
	k := make([]byte, 32)
	memguard.ScrambleBytes(k)

	// Encrypt the message.
	x, err := Encrypt(m, k)
	if err != nil {
		t.Error("expected no errors; got", err)
	}

	// Decrypt the message.
	dm := make([]byte, len(x)-Overhead)
	length, err := Decrypt(x, k, dm)
	if err != nil {
		t.Error("expected no errors; got", err)
	}
	if length != len(x)-Overhead {
		t.Error("unexpected plaintext length; got", length)
	}

	// Verify that the plaintexts match.
	if !bytes.Equal(m, dm) {
		t.Error("decrypted plaintext does not match original")
	}

	// Attempt decryption /w buffer that is too small to hold the output.
	out := make([]byte, len(x)-Overhead-1)
	length, err = Decrypt(x, k, out)
	if err != ErrBufferTooSmall {
		t.Error("expected error; got", err)
	}
	if length != 0 {
		t.Error("expected zero length; got", length)
	}

	// Construct a buffer that has the correct capacity but a smaller length.
	out = make([]byte, len(x)-Overhead)
	smallOut := out[:2]
	if len(smallOut) != 2 || cap(smallOut) != len(x)-Overhead {
		t.Error("invalid construction for test")
	}
	length, err = Decrypt(x, k, smallOut)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if length != len(x)-Overhead {
		t.Error("unexpected length; got", length)
	}
	if !bytes.Equal(m, smallOut[:len(x)-Overhead]) {
		t.Error("decrypted plaintext does not match original")
	}

	// Generate an incorrect key.
	ik := make([]byte, 32)
	memguard.ScrambleBytes(ik)

	// Attempt decryption with the incorrect key.
	length, err = Decrypt(x, ik, dm)
	if length != 0 {
		t.Error("expected length = 0; got", length)
	}
	if err != ErrDecryptionFailed {
		t.Error("expected error with incorrect key; got", err)
	}

	// Modify the ciphertext somewhat.
	for i := range x {
		if i%32 == 0 {
			x[i] = 0xdb
		}
	}

	// Attempt decryption of the invalid ciphertext with the correct key.
	length, err = Decrypt(x, k, dm)
	if length != 0 {
		t.Error("expected length = 0; got", length)
	}
	if err != ErrDecryptionFailed {
		t.Error("expected error with modified ciphertext; got", err)
	}

	// Generate a key of an invalid length.
	ik = make([]byte, 16)
	memguard.ScrambleBytes(ik)

	// Attempt encryption with the invalid key.
	ix, err := Encrypt(m, ik)
	if err != ErrInvalidKeyLength {
		t.Error("expected error with invalid key; got", err)
	}
	if ix != nil {
		t.Error("expected nil ciphertext; got", dm)
	}

	// Attempt decryption with the invalid key.
	length, err = Decrypt(x, ik, dm)
	if length != 0 {
		t.Error("expected length = 0; got", length)
	}
	if err != ErrInvalidKeyLength {
		t.Error("expected error with invalid key; got", err)
	}
}
