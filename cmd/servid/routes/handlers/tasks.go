package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tasks17-server/internal/tasks"
)

func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	tList := tasks.LoadTasks(limit)

	jsonResponse(tList, w)
}

func TaskStateListHandler(w http.ResponseWriter, r *http.Request) {
	tList := tasks.LoadTaskStateList(1)

	jsonResponse(tList, w)
}

func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t tasks.Task
	decoder.Decode(&t)

	sqlResult, err := tasks.CreateTask(t)

	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	lastInsertId, _ := sqlResult.LastInsertId()
	fmt.Fprintf(w, "{\"id\": %d}", lastInsertId)
}

func TaskCommentCreateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var comment tasks.TaskComment
	decoder.Decode(&comment)

	sqlResult, err := tasks.LeaveComment(comment)

	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	lastInsertId, _ := sqlResult.LastInsertId()
	fmt.Fprintf(w, "{\"id\": %d}", lastInsertId)
}

func TaskCommentLoadHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	vars := mux.Vars(r)

	comments := tasks.LoadComments(vars["task"], limit)

	jsonResponse(comments, w)
}
