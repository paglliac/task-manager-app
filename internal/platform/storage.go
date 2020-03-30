package platform

import (
	"database/sql"
	"fmt"
	"log"
)

type Storage struct {
	*sql.DB
}

func InitDB(dbHost string) (Storage, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:tasks17@%s/tasks17?sslmode=disable", dbHost))

	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connection with database established")

	return Storage{db}, err
}
