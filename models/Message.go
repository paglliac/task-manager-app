package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Message struct {
	Id         int
	Author     string
	Message    string
	OccurredOn time.Time
}

func LoadMessages(db *sql.DB, limit int) ([]*Message, error) {
	mList := make([]*Message, 0)

	s := fmt.Sprintf("SELECT id, author, message, occurred_on from messages LIMIT %d", limit)
	rows, err := db.Query(s)

	defer rows.Close()

	if err != nil {
		log.Panic(err)
		return mList, err
	}

	for rows.Next() {
		message := new(Message)
		rows.Scan(&message.Id, &message.Author, &message.Message, &message.OccurredOn)
		mList = append(mList, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mList, nil
}
