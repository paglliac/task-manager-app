package platform

import (
	"database/sql"
	"log"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3307)/golang?parseTime=true")

	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connection with database established")

	return db, err
}
