package tree

func (t *Tree) Count() int {
	return t.count(1)
}

func (t *Tree) count(nodeID int64) int {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		return 0
	}

	var size int
	if node.LeftChild.Valid {
		size += t.count(node.LeftChild.Int64)
	}
	if node.RightChild.Valid {
		size += t.count(node.RightChild.Int64)
	}
	if len(node.LeftObject) > 0 {
		size += 1
	}
	if len(node.RightObject) > 0 {
		size += 1
	}

	return size
}
