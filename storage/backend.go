package storage

const (
	// Configures storage backends to overwrite data with random bytes when deleting
	// instead of removing the file. This is useful for deniable encryption.
	OverwriteToDelete = true
)

const (
	MB = 1_000_000
)

type Backend interface {
	Keys() ([]string, error) // TODO: is this needed?
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
