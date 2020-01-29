package models

import (
	"ResearchGolang/internal/platform"
	"fmt"
	"log"
	"runtime"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func LoadUsers(limit int) ([]User, error) {
	userList := make([]User, 0)

	s := fmt.Sprintf("SELECT * from users LIMIT %d", limit)
	rows, err := platform.Db.Query(s)

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

	PrintMemUsage()

	return userList, nil
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v B", m.Alloc)
	fmt.Printf("\tTotalAlloc = %v B", m.TotalAlloc)
	fmt.Printf("\tSys = %v B", m.Sys)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
