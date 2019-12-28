package models

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func LoadUsers(db *sql.DB, limit int) ([]User, error) {
	userList := make([]User, 0)

	s := fmt.Sprintf("SELECT * from users LIMIT %d", limit)
	rows, err := db.Query(s)

	defer rows.Close()

	if err != nil {
		log.Panic(err)
		return userList, err
	}

	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name, &user.Email)
		userList = append(userList, user)
	}

	return userList, nil
}
