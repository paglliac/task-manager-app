package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	var err error

	db, err = sql.Open("mysql", "root@tcp(localhost:3307)/golang?parseTime=true")

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
		return nil, err
	}

	log.Println("Connection with database established")

	return db, err
}
