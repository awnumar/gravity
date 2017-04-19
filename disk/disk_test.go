package disk

import "testing"

func TestGetFileInfo(t *testing.T) {
	// Test when it exists.
	info, err := GetFileInfo(".")
	if err != nil {
		t.Error("! Unexpected error:", err)
	}

	if info == nil {
		t.Error("! File is nil")
	}

	// Test when it doesn't exist.
	info, err = GetFileInfo("test.err")
	if err == nil {
		t.Error("! Expected err; got nil")
	}

	if info != nil {
		t.Error("! Info is not nil")
	}
}
