package main

import (
	"log"
	"miru/pkg/tree"
	"os"
	"path/filepath"
)

func insert(dbName, folder string) error {
	tree, err := tree.New(dbName)
	if err != nil {
		return err
	}
	defer tree.Close()

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			if err := tree.Add(path); err != nil {
				log.Print(err)
			}
		}
		return nil
	})

	return err
}
