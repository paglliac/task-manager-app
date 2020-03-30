package auth

import (
	"tasks17-server/internal/platform"
)

var storage platform.Storage

func InitAuthModule(db platform.Storage) {
	storage = db
}

func CheckAuth(email string, password string) (id int, err error) {
	row := storage.QueryRow("SELECT id from users where email=$1 and password=$2", email, password)

	err = row.Scan(&id)

	return id, err
}
