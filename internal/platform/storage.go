package platform

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func InitDB() error {
	var err error

	Db, err = sql.Open("mysql", "root@tcp(localhost:3307)/golang?parseTime=true")

	if err != nil {
		log.Panic(err)
		return err
	}

	if err := Db.Ping(); err != nil {
		log.Panic(err)
		return err
	}

	log.Println("Connection with database established")

	return err
}
