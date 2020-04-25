package tree

import (
	"miru/pkg/serialization"
	"miru/pkg/storage"
	"sync"
)

type Tree struct {
	mu         sync.Mutex
	serializer serialization.Serializer
	storage    storage.Storage
}

// New creates a new tree
// It should be closed after being used
// It consider the root element to have ID equals 1
func New(storage storage.Storage) (*Tree, error) {
	return &Tree{
		storage:    storage,
		serializer: serialization.NewGzipSerializer(),
	}, nil
}
