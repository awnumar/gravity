package crypto

import (
	"bytes"
	"testing"
)

func TestPad(t *testing.T) {
	text := []byte("yellow submarine") // 16 bytes

	// Test when padTo < len(text)
	padded, err := Pad(text, 15)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; padded:", padded)
	}

	// Test when padTo == len(text)
	padded, err = Pad(text, 16)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; padded:", padded)
	}

	// Test when padTo-1 = len(text)
	padded, err = Pad(text, 17)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 17 {
		t.Error("expected length of padded=32; got", len(padded))
	}

	// Test when padTo > len(text)
	padded, err = Pad(text, 32)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 32 {
		t.Error("expected length of padded=32; got", len(padded))
	}

	// Test when padTo >> len(text)
	padded, err = Pad(text, 4096)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 4096 {
		t.Error("expected length of padded=32; got", len(padded))
	}
}

func TestUnpad(t *testing.T) {
	text := []byte("yellow submarine") // 16 bytes

	// Test when len(text) == padTo-1
	padded, _ := Pad(text, 17)
	unpadded, err := Unpad(padded)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !bytes.Equal(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test when len(text) < padTo
	padded, err = Pad(text, 32)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	unpadded, err = Unpad(padded)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !bytes.Equal(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test when len(text) << padTo
	padded, err = Pad(text, 4096)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	unpadded, err = Unpad(padded)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !bytes.Equal(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test invalid padding.
	unpadded, err = Unpad(text)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; unpadded:", unpadded)
	}
	if unpadded != nil {
		t.Error("Expected unpadded to be nil; unpadded =", unpadded)
	}
}
