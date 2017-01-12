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
	if err.Error() != "[!] The help command requires an argument" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "help", "add"}
	mode, err = ParseArgs(args)
	if err.Error() != "help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "help", "get"}
	mode, err = ParseArgs(args)
	if err.Error() != "help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "help", "forget"}
	mode, err = ParseArgs(args)
	if err.Error() != "help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "help", "test"}
	mode, err = ParseArgs(args)
	if err.Error() != "[!] Invalid argument to help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "add"}
	mode, err = ParseArgs(args)
	if err != nil {
		t.Error("Unexpected error; got", err)
	}
	if mode != "add" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "get"}
	mode, err = ParseArgs(args)
	if err != nil {
		t.Error("Unexpected error; got", err)
	}
	if mode != "get" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "forget"}
	mode, err = ParseArgs(args)
	if err != nil {
		t.Error("Unexpected error; got", err)
	}
	if mode != "forget" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "test"}
	mode, err = ParseArgs(args)
	if err.Error() != "[!] Invalid argument" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}
}
