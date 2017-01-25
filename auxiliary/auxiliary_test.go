package auxiliary

import "testing"

func TestParseArgs(t *testing.T) {
	args := []string{"./pocket"}
	mode, _, err := ParseArgs(args)
	if err.Error() != "help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "help"}
	mode, _, err = ParseArgs(args)
	if err.Error() != "help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args2 := []string{"add", "get", "forget"}
	for _, arg := range args2 {
		args = []string{"./pocket", "help", arg}
		mode, _, err = ParseArgs(args)
		if err.Error() != "help" {
			t.Error("Expected error; got ", err)
		}
		if mode != "" {
			t.Error("Expected empty mode; got", mode)
		}

		args = []string{"./pocket", arg}
		mode, _, err = ParseArgs(args)
		if err != nil {
			t.Error("Unexpected error; got", err)
		}
		if mode != arg {
			t.Error("Expected empty mode; got", mode)
		}
	}

	args = []string{"./pocket", "help", "test"}
	mode, _, err = ParseArgs(args)
	if err.Error() != "[!] Invalid mode passed to help" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "test"}
	mode, _, err = ParseArgs(args)
	if err.Error() != "[!] Invalid option" {
		t.Error("Expected error; got ", err)
	}
	if mode != "" {
		t.Error("Expected empty mode; got", mode)
	}

	args = []string{"./pocket", "get", "-c"}
	_, _, err = ParseArgs(args)
	if err.Error() != "[!] Nothing passed to -c" {
		t.Error("Expected error; got ", err)
	}

	args = []string{"./pocket", "get", "-c", "N,r,p,test"}
	_, _, err = ParseArgs(args)
	if err.Error() != "[!] Invalid number of arguments passed to -c" {
		t.Error("Expected error; got ", err)
	}

	args2 = []string{"N,8,1", "18,r,1", "18,8,p"}
	for _, arg := range args2 {
		args = []string{"./pocket", "get", "-c", arg}
		_, _, err = ParseArgs(args)
		if err.Error() != "[!] Arguments to -c must be integers" {
			t.Error("Expected error; got ", err)
		}
	}

	args = []string{"./pocket", "get", "-c", "1,8,1"}
	_, _, err = ParseArgs(args)
	if err.Error() != "[!] N must be more than 1" {
		t.Error("Expected error; got ", err)
	}

	args = []string{"./pocket", "get", "-c", "18,8,1"}
	_, costFactor, err := ParseArgs(args)
	if err != nil {
		t.Error("Unexpected error; got", err)
	}
	if costFactor["N"] != 18 || costFactor["r"] != 8 || costFactor["p"] != 1 {
		t.Error("Expected params to be 18,8,1; got", costFactor)
	}
}
