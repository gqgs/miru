package main

import (
	"miru/pkg/database"
)

func insert(dbName, folder string) error {
	db, err := database.Open(dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
