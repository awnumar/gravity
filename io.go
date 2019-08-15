package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/awnumar/memguard"
)

func input(prompt string) *memguard.LockedBuffer {
	fmt.Printf(prompt)
	return memguard.NewBufferFromReaderUntil(os.Stdin, '\n')
}

func prompt() string {
	stdin := bufio.NewReader(os.Stdin)
	i, err := stdin.ReadString('\n')
	if err != nil {
		memguard.SafePanic(err)
	}
	return i[:len(i)-1]
}
