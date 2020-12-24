package tree

import (
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
)

func (t *Tree) Migrate() error {
	return t.migrate(1)
}

func (t *Tree) migrate(nodeID int64) error {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		return err
	}
	if node.LeftObject != nil {
		left, err := image.Deserialize(*node.LeftObject)
		if err != nil {
			return err
		}
		if err := t.storage.SetObject(nodeID, storage.Left, left); err != nil {
			return err
		}
	}
	if node.RightObject != nil {
		right, err := image.Deserialize(*node.RightObject)
		if err != nil {
			return err
		}
		if err := t.storage.SetObject(nodeID, storage.Right, right); err != nil {
			return err
		}
	}

	if node.LeftChild != nil {
		if err := t.migrate(*node.LeftChild); err != nil {
			return err
		}
	}

	if node.RightChild != nil {
		if err := t.migrate(*node.RightChild); err != nil {
			return err
		}
	}

	return nil
}
