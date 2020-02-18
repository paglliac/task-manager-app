package platform

import (
	"database/sql"
	"log"
)

type Storage struct {
	*sql.DB
}

func InitDB() (Storage, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3307)/golang?parseTime=true")

	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connection with database established")

	return Storage{db}, err
}
