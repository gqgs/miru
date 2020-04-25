package main

import (
	"fmt"
	"miru/pkg/image"
	"miru/pkg/tree"
	"os/exec"
)

func search(dbName, filename string, accuracy, limit uint, open bool) error {
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

	top := res.Top(limit)
	for _, t := range top {
		fmt.Println(t)
	}

	if open && len(top) > 0 {
		return exec.Command("xdg-open", top[0].Filename).Start()
	}

	return nil
}
