package storage

import (
	"encoding"
	"errors"
)

var ErrIsEmpty = errors.New("pkg/storage: storage is empty")

type Position int

const (
	Left Position = iota
	Right
)

type Node struct {
	LeftObject, RightObject *[]byte
	LeftChild, RightChild   *int64
}

type Storage interface {
	Get(nodeID int64) (node *Node, err error)
	SetObject(nodeID int64, position Position, marshaler encoding.BinaryMarshaler) error
	SetChild(nodeID int64, position Position, child int64) error
	NewNode(marshaler encoding.BinaryMarshaler) (nodeID int64, err error)
}
