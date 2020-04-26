package tree

import (
	"encoding"
	"sync"

	"github.com/gqgs/miru/pkg/storage"
)

type Tree struct {
	mu      sync.Mutex
	storage storage.Storage
}

// Compare receives an encoded version of the element stored in the tree
// `distance` controls how the elements will be partitioned in the hyperplane
// Ideally, the method should follow the requirements of a distance function
// See: https://en.wikipedia.org/wiki/Metric_(mathematics)#Definition
type Comparer interface {
	encoding.BinaryMarshaler
	Compare(element []byte) (distance float64, comparedElement string, err error)
}

// New creates a new tree
// It consider the root element to have ID equals 1
func New(storage storage.Storage) *Tree {
	return &Tree{
		storage: storage,
	}
}
