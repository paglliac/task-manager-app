package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
	authorId, _ := strconv.Atoi(r.Header.Get("Authorization"))
	tList := tasks.LoadTaskStateList(authorId)

	jsonResponse(tList, w)
}

func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t tasks.Task
	decoder.Decode(&t)

	authorId, _ := strconv.Atoi(r.Header.Get("Authorization"))
	t.AuthorId = authorId

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
	authorId, err := strconv.Atoi(r.Header.Get("Authorization"))
	comment.Author = authorId

	sqlResult, err := tasks.LeaveComment(comment)

	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	lastInsertId, _ := sqlResult.LastInsertId()
	fmt.Fprintf(w, "{\"id\": %d}", lastInsertId)
}

func TaskLoadHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Task     tasks.Task          `json:"task"`
		Comments []tasks.TaskComment `json:"comments"`
	}

	var rs response

	vars := mux.Vars(r)

	taskId := vars["task"]
	rs.Comments = tasks.LoadComments(taskId)
	rs.Task = tasks.LoadTask(taskId)

	jsonResponse(rs, w)
}

func TaskUpdateLastCommentHandler(w http.ResponseWriter, r *http.Request) {
	type updateLastComment struct {
		TaskId    string `json:"task_id"`
		CommentId string `json:"comment_id"`
		UserId    int    `json:"-"`
	}

	userId, _ := strconv.Atoi(r.Header.Get("Authorization"))

	decoder := json.NewDecoder(r.Body)
	var t updateLastComment
	decoder.Decode(&t)
	t.UserId = userId

	tasks.UpdateLastWatchedComment(t.UserId, t.TaskId, t.CommentId)
}

func TaskUpdateDescription(w http.ResponseWriter, r *http.Request) {
	type newDescription struct {
		Description string `json:"description"`
	}

	d := json.NewDecoder(r.Body)

	var nd newDescription
	d.Decode(&nd)
	vars := mux.Vars(r)

	taskId := vars["task"]

	task := tasks.Task{Id: taskId}
	err := task.UpdateDescription(nd.Description)

	if err != nil {
		log.Println(err)
	}
}

func TaskClose(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId := vars["task"]

	task := tasks.Task{Id: taskId}
	err := task.Close()

	if err != nil {
		log.Println(err)
	}
}
