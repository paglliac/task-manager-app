package main

import (
	"ResearchGolang/models"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

func LoggedHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		handler(w, r)
	})
}

func main() {
	var err error

	db, err := InitDB()

	if err != nil {
		panic(err)
	}

	LoggedHandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", 404)
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil {
			limit = 10
		}

		userList, err := models.LoadUsers(db, limit)

		if err != nil {
			http.Error(w, "Unexpected error", 500)
			return
		}

		jsonResponse, _ := json.Marshal(userList)

		fmt.Fprintf(w, string(jsonResponse))
	})

	LoggedHandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			http.Error(w, "Not found", 404)
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil {
			limit = 10
		}

		messages, err := models.LoadMessages(db, limit)

		if err != nil {
			http.Error(w, "Unexpected error", 500)
			return
		}

		jsonResponse, _ := json.Marshal(messages)

		fmt.Fprintf(w, string(jsonResponse))
	})

	LoggedHandleFunc("/messages/add", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages/add" {
			http.Error(w, "Not found", 404)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var m models.Message
		decoder.Decode(&m)

		if m.OccurredOn.IsZero() {
			http.Error(w, "Bad Request", 400)
			return
		}

		sqlResult, err := models.SaveMessage(db, m)

		if err != nil {
			http.Error(w, "Something went wrong", 500)
			return
		}

		lastInsertId, _ := sqlResult.LastInsertId()
		fmt.Fprintf(w, "{\"id\": %d}", lastInsertId)

	})

	log.Println("Server have been started listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
