package auth

import (
	"log"
	"tasks17-server/internal/platform"
)

var storage platform.Storage

func InitAuthModule(db platform.Storage) {
	storage = db
}

func CheckAuth(email string, password string) (id int, err error) {
	log.Printf("Email is %s, password is %s", email, password)
	row := storage.QueryRow("SELECT id from users where email=? and password=?", email, password)

	err = row.Scan(&id)

	return id, err
}
