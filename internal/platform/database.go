package platform

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDb(dbHost string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:tasks17@%s/tasks17?sslmode=disable", dbHost))

	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connection with database established")

	return db, err
}
