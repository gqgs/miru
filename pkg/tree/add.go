package tree

import (
	"github.com/gqgs/miru/pkg/storage"
)

// Add recursively traversals the tree to find the
// correct insert position for the element
func (t *Tree) Add(comparer Comparer) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.add(1, comparer)
}

func (t *Tree) add(nodeID int64, comparer Comparer) error {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		if err == storage.ErrIsEmpty {
			_, err = t.storage.NewNode(comparer)
			return err
		}
		return err
	}

	var cmp0 float64
	if len(node.LeftObject) == 0 {
		return t.storage.SetObject(nodeID, storage.Left, comparer)
	}
	if cmp0, _, err = comparer.Compare(node.LeftObject); err != nil {
		return err
	}
	if cmp0 == 0 {
		return t.storage.SetObject(nodeID, storage.Left, comparer)
	}

	var cmp1 float64
	if len(node.RightObject) == 0 {
		return t.storage.SetObject(nodeID, storage.Right, comparer)
	}
	if cmp1, _, err = comparer.Compare(node.RightObject); err != nil {
		return err
	}
	if cmp1 == 0 {
		return t.storage.SetObject(nodeID, storage.Right, comparer)
	}

	if cmp0 < cmp1 {
		if !node.LeftChild.Valid {
			lastID, err := t.storage.NewNode(comparer)
			if err != nil {
				return err
			}
			return t.storage.SetChild(nodeID, storage.Left, lastID)
		}
		return t.add(node.LeftChild.Int64, comparer)
	}
	if !node.RightChild.Valid {
		lastID, err := t.storage.NewNode(comparer)
		if err != nil {
			return err
		}
		return t.storage.SetChild(nodeID, storage.Right, lastID)
	}
	return t.add(node.RightChild.Int64, comparer)
}
