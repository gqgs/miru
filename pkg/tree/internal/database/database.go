package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tree (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			image0 BLOB NULL,
			image1 BLOB NULL,
			left INTEGER NULL,
			right INTEGER NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
