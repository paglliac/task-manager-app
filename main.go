package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

type User struct {
	Id    int
	Name  string
	Email string
}

func main() {
	var err error

	db, err = sql.Open("mysql", "root@tcp(localhost:3307)/golang")

	if err != nil {
		log.Panic(err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
		return
	}

	log.Println("Connection with database established")

	rows, err := db.Query("SELECT * from users")
	defer rows.Close()

	if err != nil {
		log.Panic(err)
		return
	}

	userList := make([]User, 0)

	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name, &user.Email)
		log.Println(user.Name)
		userList = append(userList, user)
	}

	jsonResponse, _ := json.Marshal(userList)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)

		if r.URL.Path != "/" {
			http.Error(w, "Not found", 404)
		}

		fmt.Fprintf(w, string(jsonResponse))
	})

	http.ListenAndServe(":8080", nil)
}
