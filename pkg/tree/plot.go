package tree

import (
	"fmt"
	"io"
)

// Writes to the writer a graph following the DOT spec
// https://en.wikipedia.org/wiki/DOT_(graph_description_language)
func (t *Tree) Plot(writer io.Writer) error {
	t.writer = writer
	if _, err := io.WriteString(t.writer, "digraph {\n"); err != nil {
		return err
	}
	if err := t.plot("", 1); err != nil {
		return err
	}
	if _, err := io.WriteString(t.writer, "}\n"); err != nil {
		return err
	}
	return nil
}

func (t *Tree) plot(parent string, nodeID int64) error {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		return err
	}
	var objects uint
	if len(node.LeftObject) > 0 {
		objects++
	}
	if len(node.RightObject) > 0 {
		objects++
	}

	name := fmt.Sprintf("node_%d_%d", nodeID, objects)

	// Define node
	fmt.Fprintf(t.writer, "%s [shape = box];\n", name)

	// Define edge
	if parent != "" {
		fmt.Fprintf(t.writer, "%s -> %s;\n", parent, name)
	}

	if node.LeftChild.Valid {
		if err = t.plot(name, node.LeftChild.Int64); err != nil {
			return err
		}
	}
	if node.RightChild.Valid {
		if err = t.plot(name, node.RightChild.Int64); err != nil {
			return err
		}
	}

	return nil
}
