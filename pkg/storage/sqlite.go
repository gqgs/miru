package storage

import (
	"database/sql"
	"encoding"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	db         *sql.DB
	stmt       *sql.Stmt
	compressor compress.Compressor
	cache      cache.Cache
}

// Should be closed after being used
func NewSqliteStorage(dbName string, compressor compress.Compressor, cache cache.Cache) (*sqliteStorage, error) {
	db, err := sql.Open("sqlite3", dbName+"?_synchronous=off&_journal_mode=off&cache=shared")
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tree (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			image0 BLOB NULL,
			image1 BLOB NULL,
			left INTEGER NULL,
			right INTEGER NULL
		)
	`); err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(
		`SELECT image0, image1, left, right
		FROM tree
		WHERE id = ?
		`)
	if err != nil {
		return nil, err
	}

	return &sqliteStorage{
		db:         db,
		stmt:       stmt,
		compressor: compressor,
		cache:      cache,
	}, nil
}

func (s *sqliteStorage) Close() error {
	_ = s.stmt.Close()
	return s.db.Close()
}

func (s *sqliteStorage) Get(nodeID int64) (*Node, error) {
	value, ok := s.cache.Get(nodeID)
	if ok {
		return value.(*Node), nil
	}

	node := new(Node)
	err := s.stmt.QueryRow(nodeID).
		Scan(&node.LeftObject, &node.RightObject, &node.LeftChild, &node.RightChild)
	if err != nil {
		if err == sql.ErrNoRows && nodeID == 1 {
			return nil, ErrIsEmpty
		}
		return nil, err
	}
	if node.LeftObject != nil {
		decompressed, err := s.compressor.Decompress(*node.LeftObject)
		if err != nil {
			return nil, err
		}
		*node.LeftObject = decompressed
	}
	if node.RightObject != nil {
		decompressed, err := s.compressor.Decompress(*node.RightObject)
		if err != nil {
			return nil, err
		}
		*node.RightObject = decompressed
	}

	s.cache.Add(nodeID, node)
	return node, nil
}

func (s *sqliteStorage) SetObject(nodeID int64, position Position, marshaler encoding.BinaryMarshaler) (err error) {
	s.cache.Remove(nodeID)
	b, err := marshaler.MarshalBinary()
	if err != nil {
		return err
	}
	data, err := s.compressor.Compress(b)
	if err != nil {
		return err
	}

	switch position {
	case Left:
		_, err = s.db.Exec(`UPDATE tree SET image0 = ? WHERE id = ?`, data, nodeID)
	case Right:
		_, err = s.db.Exec(`UPDATE tree SET image1 = ? WHERE id = ?`, data, nodeID)
	}
	return
}

func (s *sqliteStorage) SetChild(nodeID int64, position Position, child int64) (err error) {
	s.cache.Remove(nodeID)
	switch position {
	case Left:
		_, err = s.db.Exec("UPDATE tree SET left = ? WHERE id = ?", child, nodeID)
	case Right:
		_, err = s.db.Exec("UPDATE tree SET right = ? WHERE id = ?", child, nodeID)
	}
	return
}

func (s *sqliteStorage) NewNode(marshaler encoding.BinaryMarshaler) (nodeID int64, err error) {
	b, err := marshaler.MarshalBinary()
	if err != nil {
		return 0, err
	}
	data, err := s.compressor.Compress(b)
	if err != nil {
		return 0, err
	}

	result, err := s.db.Exec(`INSERT INTO tree (image0) VALUES (?)`, data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
