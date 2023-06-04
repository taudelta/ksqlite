package ksqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) readAll() ([]string, error) {
	db, err := sql.Open("sqlite3", "reader.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name FROM reader_db ORDER BY id")
	if err != nil {
		return nil, err
	}

	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}

	return result, nil
}

func (r *Reader) Run() error {
	topic := "my-topic"
	partition := 0

	fmt.Println("start to send")

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	messages := []kafka.Message{}
	result, err := r.readAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result {
		messages = append(messages, kafka.Message{
			Value: []byte(r),
		})
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	fmt.Println("send ok")

	return nil
}
