package tree

import (
	"miru/pkg/image"
	"miru/pkg/storage"
)

// Add recursively traversals the tree to find the
// correct insert position for the image
func (t *Tree) Add(img *image.Image) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.add(1, img)
}

func (t *Tree) add(nodeID int64, img *image.Image) error {
	node, err := t.storage.Get(nodeID)
	switch err {
	case storage.ErrIsEmpty:
		data, err := t.serializer.Marshal(img)
		if err != nil {
			return err
		}
		_, err = t.storage.NewNode(data)
		return err
	case nil:
		if node.LeftObject == nil {
			data, err := t.serializer.Marshal(img)
			if err != nil {
				return err
			}
			err = t.storage.UpdateObject(nodeID, storage.Left, data)
			return err
		}
		if node.RightObject == nil {
			data, err := t.serializer.Marshal(img)
			if err != nil {
				return err
			}
			err = t.storage.UpdateObject(nodeID, storage.Right, data)
			return err
		}
		var dbImage0, dbImage1 image.Image
		if err = t.serializer.Unmarshal(*node.LeftObject, &dbImage0); err != nil {
			return err
		}
		if err = t.serializer.Unmarshal(*node.RightObject, &dbImage1); err != nil {
			return err
		}
		cmp0 := image.Compare(img, &dbImage0)
		cmp1 := image.Compare(img, &dbImage1)
		if cmp0 < cmp1 {
			if node.LeftChild == nil {
				data, err := t.serializer.Marshal(img)
				if err != nil {
					return err
				}
				lastID, err := t.storage.NewNode(data)
				if err != nil {
					return err
				}
				err = t.storage.UpdateChild(nodeID, storage.Left, lastID)
				return err
			}
			return t.add(*node.LeftChild, img)
		}
		if node.RightChild == nil {
			data, err := t.serializer.Marshal(img)
			if err != nil {
				return err
			}
			lastID, err := t.storage.NewNode(data)
			if err != nil {
				return err
			}
			err = t.storage.UpdateChild(nodeID, storage.Right, lastID)
			return err
		}
		return t.add(*node.RightChild, img)
	}
	return err
}
