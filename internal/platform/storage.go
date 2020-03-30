package platform

import (
	"database/sql"
	"log"
)

type Storage struct {
	*sql.DB
}

func InitDB() (Storage, error) {
	connStr := "postgres://postgres:tasks17@localhost/tasks17?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	//db, err := sql.Open("postgres", "root:tasks17@tcp(localhost:3308)/tasks17?parseTime=true")

	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connection with database established")

	return Storage{db}, err
}
