package ksqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func createWriterTestDB() {
	db, err := sql.Open("sqlite3", "writer.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS writer_db (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name string
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM writer_db")
	if err != nil {
		panic(err)
	}
}
func TestWriter(t *testing.T) {
	createWriterTestDB()

	writer := NewWriter()
	writer.Run(2)

	db, err := sql.Open("sqlite3", "writer.db")
	if err != nil {
		t.Error(err)
	}

	rows, err := db.Query("SELECT name FROM writer_db ORDER BY name")
	if err != nil {
		t.Error(err)
	}

	expected := []string{
		"test1",
		"test2",
	}

	i := 0
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			t.Error(err)
		}

		if row != expected[i] {
			t.Errorf("row must be equals, got: %s, expect: %s", row, expected[i])
			return
		}

		i++
	}
}
