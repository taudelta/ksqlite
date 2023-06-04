package ksqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func createTestDB() {
	db, err := sql.Open("sqlite3", "reader.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reader_db (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name string
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM reader_db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO reader_db (name) VALUES (?)", "test1")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO reader_db (name) VALUES (?)", "test2")
	if err != nil {
		panic(err)
	}
}

func TestReader(t *testing.T) {
	createTestDB()

	reader := NewReader()
	reader.Run()
}
