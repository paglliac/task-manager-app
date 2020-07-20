package platform

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDb(config DbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Db))

	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connection with database established")

	return db, err
}
