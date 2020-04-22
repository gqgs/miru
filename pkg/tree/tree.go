package tree

import (
	"database/sql"
	"miru/pkg/tree/internal/database"
)

type Tree struct {
	db   *sql.DB
	stmt *sql.Stmt
}

// New creates a new tree
// It should be closed after being used
func New(dbName string) (*Tree, error) {
	db, err := database.Open(dbName)
	if err != nil {
		return nil, err
	}
	return &Tree{
		db: db,
	}, nil
}

func (t *Tree) Close() error {
	return t.db.Close()
}
