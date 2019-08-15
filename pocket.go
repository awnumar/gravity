package main

import (
	"unsafe"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"

	"github.com/awnumar/memguard"
)

// Argon2id parameters
const (
	iters   = 128       // 128 iterations
	memory  = 64 * 1024 // over 64 MiB of memory
	threads = 4         // by 4 threads
)

// Pocket defines a folder within which data can be stored. A particular folder is uniquely identified by a key.
type Pocket struct {
	ID  *memguard.Enclave
	Key *memguard.Enclave
}

// GetPocket takes a key and derives a unique folder within which data may be stored.
func GetPocket(key *memguard.LockedBuffer) *Pocket {
	root := memguard.NewBufferFromBytes(argon2.IDKey(key.Bytes(), []byte{}, iters, memory, threads, 64))
	go key.Destroy()
	defer root.Destroy()
	root.Melt()
	return &Pocket{memguard.NewEnclave(root.Bytes()[:32]), memguard.NewEnclave(root.Bytes()[32:])}
}

// Identifier specifies the values used to derive identifiers.
type Identifier struct {
	root  [32]byte
	file  uint64
	chunk uint64
}

// Identifier initialises a new Identifier struct within some securely allocated memory.
func (p *Pocket) Identifier() (*Identifier, *memguard.LockedBuffer, error) {
	// Initialise an Identifier struct within some securely allocated memory.
	i := new(Identifier)
	b := memguard.NewBuffer(int(unsafe.Sizeof(*i)))
	i = (*Identifier)(unsafe.Pointer(&b.Bytes()[0]))

	// Unlock the root id and copy it into the structure.
	root, err := p.ID.Open()
	if err != nil {
		return nil, nil, err
	}
	copy(i.root[:], root.Bytes())
	go root.Destroy()

	return i, b, nil
}

// Derive derives the identifier for a given (key, file, chunk) pair.
func (i *Identifier) Derive(memory *memguard.LockedBuffer, file, chunk uint64) []byte {
	i.file = file
	i.chunk = chunk
	id := blake2b.Sum256(memory.Bytes())
	return id[:]
}
