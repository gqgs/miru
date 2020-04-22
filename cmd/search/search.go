package main

import (
	"container/heap"
	"fmt"
	"miru/pkg/tree"
)

func search(dbName, fileName string, accuracy, limit int) error {
	tree, err := tree.New(dbName)
	if err != nil {
		return err
	}
	defer tree.Close()

	res, err := tree.Search(fileName, accuracy)
	if err != nil {
		return err
	}

	for heap.Init(&res); len(res) > 0 && limit > 0; limit-- {
		fmt.Println(heap.Pop(&res))
	}
	return nil
}
