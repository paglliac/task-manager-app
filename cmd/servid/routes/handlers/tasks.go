package handlers

import (
	"ResearchGolang/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasks" {
		http.Error(w, "Not found", 404)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	tasks := models.LoadTasks(limit)

	jsonResponse, _ := json.Marshal(tasks)

	_, err = fmt.Fprintf(w, string(jsonResponse))

	if err != nil {
		log.Println("Can't write response in response writer", err)
	}
}
