package main

import (
	"fmt"
	"miru/pkg/image"
	"miru/pkg/tree"
)

func search(dbName, filename string, accuracy, limit uint) error {
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

	for _, top := range res.Top(limit) {
		fmt.Println(top)
	}
	return nil
}
