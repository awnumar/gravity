package store

// Store defines an interface of which all stores must implement to be used
// within Pocket.
type Store interface {
	Save(identifier, ciphertext []byte) error
	Retrieve(identifier []byte) ([]byte, error)
	Delete(identifier []byte)
	Close() error
}
