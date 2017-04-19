package crypto

import "testing"

func TestGenDecoy(t *testing.T) {
	id, ct := GenDecoy()

	// Check if they're the right length.

	if len(id) != 32 {
		t.Error("! Derived identifier incorrect length:", len(id))
	}

	if len(ct) != 4136 {
		t.Error("! Ciphertext incorrect length:", len(ct))
	}

	// Check if they're null.

	if id == nil {
		t.Error("Derived identifier is nil.")
	}

	if ct == nil {
		t.Error("Ciphertext is nil.")
	}
}
