package main

import (
	"container/heap"
	"fmt"
	"miru/pkg/image"
	"miru/pkg/tree"
)

func search(dbName, filename string, accuracy, limit int) error {
	tree, err := tree.New(dbName)
	if err != nil {
		return err
	}
	defer tree.Close()

	img, err := image.Load(filename)
	if err != nil {
		return err
	}

	res, err := tree.Search(img, accuracy)
	if err != nil {
		return err
	}

	for heap.Init(&res); len(res) > 0 && limit > 0; limit-- {
		fmt.Println(heap.Pop(&res))
	}
	return nil
}
