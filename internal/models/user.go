package models

import (
	"fmt"
	"log"
	"tasks17-server/internal/platform"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func LoadUsers(limit int) ([]User, error) {
	userList := make([]User, 0)

	s := fmt.Sprintf("SELECT id, name, email from users LIMIT %d", limit)
	rows, err := platform.Db.Query(s)

	defer rows.Close()

	if err != nil {
		log.Panic(err)
		return userList, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return userList, err
		}
		userList = append(userList, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}
