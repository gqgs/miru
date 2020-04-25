package tree

import (
	"miru/pkg/storage"
	"sync"
)

type Tree struct {
	mu      sync.Mutex
	storage storage.Storage
}

// Comparer receives an encoded version of the element stored in the tree
// `distance` controls how the elements will be partitioned in the hyperplane
type Comparer interface {
	Compare(element []byte) (distance float64, comparedElement string, err error)
}

// New creates a new tree
// It consider the root element to have ID equals 1
func New(storage storage.Storage) (*Tree, error) {
	return &Tree{
		storage: storage,
	}, nil
}
