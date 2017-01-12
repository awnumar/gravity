package auxiliary

import "testing"

func TestParseArgs(t *testing.T) {
	args := []string{"./pocket"}
	mode, err := ParseArgs(args)
	if err.Error() != "help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "help"}
	mode, err = ParseArgs(args)
	if err.Error() != "[!] Invalid arguments" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args2 := []string{"add", "get", "forget"}
	for _, arg := range args2 {
		args = []string{"./pocket", "help", arg}
		mode, err = ParseArgs(args)
		if err.Error() != "help" {
			t.Error("Expected error; got ", err)
		}
		if mode != "" {
			t.Error("Expected empty mode; got", mode)
		}
	}

	args = []string{"./pocket", "help", "test"}
	mode, err = ParseArgs(args)
	if err.Error() != "[!] Invalid arguments" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args2 = []string{"add", "get", "forget"}
	for _, arg := range args2 {
		args = []string{"./pocket", arg}
		mode, err = ParseArgs(args)
		if err != nil {
			t.Error("Unexpected error; got", err)
		}
		if mode != arg {
			t.Error("Expected empty mode; got", mode)
		}
	}

	args = []string{"./pocket", "test"}
	mode, err = ParseArgs(args)
	if err.Error() != "[!] Invalid arguments" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}
}
