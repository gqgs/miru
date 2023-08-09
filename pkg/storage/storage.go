package storage

import (
	"database/sql"
	"encoding"
	"errors"
	"io"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
)

var ErrIsEmpty = errors.New("pkg/storage: storage is empty")

type Position int

const (
	Left Position = iota
	Right
)

func (p Position) Object() string {
	if p == Left {
		return "image0"
	}
	return "image1"
}

func (p Position) Child() string {
	if p == Left {
		return "left"
	}
	return "right"
}

type nullInt64 struct {
	sql.NullInt64
}

type Node struct {
	LeftObject  []byte    `redis:"image0"`
	RightObject []byte    `redis:"image1"`
	LeftChild   nullInt64 `redis:"left"`
	RightChild  nullInt64 `redis:"right"`
}

type Storage interface {
	io.Closer
	Get(nodeID int64) (node *Node, err error)
	SetObject(nodeID int64, position Position, marshaler encoding.BinaryMarshaler) error
	SetChild(nodeID int64, position Position, child int64) error
	NewNode(marshaler encoding.BinaryMarshaler) (nodeID int64, err error)
}

func NewStorage(name string, dbName string, compressor compress.Compressor, cache cache.Cache) (Storage, error) {
	switch name {
	case "redis":
		return newRedisStorage(dbName, compressor, cache)
	case "sqlite":
		return newSqliteStorage(dbName, compressor, cache)
	}
	return nil, errors.New("invalid storage")
}
