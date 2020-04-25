package storage

import "errors"

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
	Get(nodeID int64) (node Node, err error)
	UpdateObject(nodeID int64, position Position, data interface{}) error
	UpdateChild(nodeID int64, position Position, child int64) error
	NewNode(data interface{}) (nodeID int64, err error)
}
