package handlers

import (
	"ResearchGolang/internal/models"
	"ResearchGolang/internal/platform"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func MessageListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/messages" {
		http.Error(w, "Not found", 404)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	messages, err := models.LoadMessages(limit)

	if err != nil {
		http.Error(w, "Unexpected error", 500)
		return
	}

	jsonResponse, _ := json.Marshal(messages)

	_, err = fmt.Fprintf(w, string(jsonResponse))

	if err != nil {
		log.Println("Can't write response in response writer", err)
	}
}

func MessageCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

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

	sqlResult, err := models.SaveMessage(m)

	b, err := json.Marshal(m)
	// todo process error

	platform.CurrentHub.Broadcast <- b

	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	lastInsertId, _ := sqlResult.LastInsertId()
	fmt.Fprintf(w, "{\"id\": %d}", lastInsertId)

}
