package storage

import (
	"database/sql"
	"encoding/json"
	"miru/pkg/compress"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	db         *sql.DB
	stmt       *sql.Stmt
	compressor compress.Compressor
}

// Should be closed after being used
func NewSqliteStorage(dbName string, compressor compress.Compressor) (*sqliteStorage, error) {
	db, err := sql.Open("sqlite3", dbName)
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
	if _, err = db.Exec(`PRAGMA synchronous = OFF`); err != nil {
		return nil, err
	}
	if _, err = db.Exec(`PRAGMA journal_mode = OFF`); err != nil {
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
	}, nil
}

func (s *sqliteStorage) Close() error {
	_ = s.stmt.Close()
	return s.db.Close()
}

func (s *sqliteStorage) Get(nodeID int64) (*Node, error) {
	node := new(Node)
	err := s.stmt.QueryRow(nodeID).
		Scan(&node.LeftObject, &node.RightObject, &node.LeftChild, &node.RightChild)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrIsEmpty
		}
		return nil, err
	}
	if node.LeftObject != nil {
		decompressed, err := s.compressor.Decompress(*node.LeftObject)
		if err != nil {
			return nil, err
		}
		*node.LeftObject = append(decompressed[0:0], decompressed...)
	}
	if node.RightObject != nil {
		decompressed, err := s.compressor.Decompress(*node.RightObject)
		if err != nil {
			return nil, err
		}
		*node.RightObject = append(decompressed[0:0], decompressed...)
	}
	return node, nil
}

func (s *sqliteStorage) SetObject(nodeID int64, position Position, obj interface{}) (err error) {
	b, err := json.Marshal(obj)
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
	switch position {
	case Left:
		_, err = s.db.Exec("UPDATE tree SET left = ? WHERE id = ?", child, nodeID)
	case Right:
		_, err = s.db.Exec("UPDATE tree SET right = ? WHERE id = ?", child, nodeID)
	}
	return
}

func (s *sqliteStorage) NewNode(obj interface{}) (nodeID int64, err error) {
	b, err := json.Marshal(obj)
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
