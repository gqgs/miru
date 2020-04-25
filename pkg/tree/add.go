package tree

import (
	"github.com/gqgs/miru/pkg/storage"
)

// Add recursively traversals the tree to find the
// correct insert position for the image
func (t *Tree) Add(comparer Comparer) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.add(1, comparer)
}

func (t *Tree) add(nodeID int64, comparer Comparer) error {
	node, err := t.storage.Get(nodeID)
	switch err {
	case storage.ErrIsEmpty:
		_, err = t.storage.NewNode(comparer)
		return err
	case nil:
		if node.LeftObject == nil {
			err = t.storage.SetObject(nodeID, storage.Left, comparer)
			return err
		}
		if node.RightObject == nil {
			err = t.storage.SetObject(nodeID, storage.Right, comparer)
			return err
		}
		var cmp0, cmp1 float64
		if cmp0, _, err = comparer.Compare(*node.LeftObject); err != nil {
			return err
		}
		if cmp1, _, err = comparer.Compare(*node.RightObject); err != nil {
			return err
		}
		if cmp0 < cmp1 {
			if node.LeftChild == nil {
				lastID, err := t.storage.NewNode(comparer)
				if err != nil {
					return err
				}
				err = t.storage.SetChild(nodeID, storage.Left, lastID)
				return err
			}
			return t.add(*node.LeftChild, comparer)
		}
		if node.RightChild == nil {
			lastID, err := t.storage.NewNode(comparer)
			if err != nil {
				return err
			}
			err = t.storage.SetChild(nodeID, storage.Right, lastID)
			return err
		}
		return t.add(*node.RightChild, comparer)
	}
	return err
}
