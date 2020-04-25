package main

import (
	"fmt"
	"miru/pkg/image"
	"miru/pkg/tree"
	"os/exec"
)

func search(o options) error {
	tree, err := tree.New(o.db)
	if err != nil {
		return err
	}
	defer tree.Close()

	img, err := image.Load(o.file)
	if err != nil {
		return err
	}

	res, err := tree.Search(img, o.accuracy)
	if err != nil {
		return err
	}

	top := res.Top(o.limit)
	for _, t := range top {
		fmt.Println(t)
	}

	if o.open && len(top) > 0 {
		return exec.Command("xdg-open", top[0].Filename).Start()
	}

	return nil
}
