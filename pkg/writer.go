package ksqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) writeAll(record string) error {
	db, err := sql.Open("sqlite3", "writer.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO writer_db (name) VALUES (?)", record)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Run(maxMessages int) {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message

	messagesCount := 0

	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))

		if err := w.writeAll(string(b[:n])); err != nil {
			log.Fatal("failed to write: ", err)
		}

		messagesCount++
		if messagesCount == maxMessages {
			break
		}
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
