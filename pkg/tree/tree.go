package tree

import (
	"database/sql"
	"miru/pkg/serialize"
	"miru/pkg/tree/internal/database"
	"sync"
)

type Tree struct {
	mu         sync.Mutex
	db         *sql.DB
	stmt       *sql.Stmt
	serializer serialize.Serializer
}

// New creates a new tree
// It should be closed after being used
func New(dbName string) (*Tree, error) {
	db, err := database.Open(dbName)
	if err != nil {
		return nil, err
	}
	return &Tree{
		db:         db,
		serializer: serialize.NewGzip(),
	}, nil
}

func (t *Tree) Close() error {
	return t.db.Close()
}
