package handlers

import (
	"ResearchGolang/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	userList, err := models.LoadUsers(limit)

	if err != nil {
		http.Error(w, "Unexpected error", 500)
		return
	}

	jsonResponse, _ := json.Marshal(userList)

	_, err = fmt.Fprintf(w, string(jsonResponse))

	if err != nil {
		log.Println("Can't write response in response writer", err)
	}

}
