package tasks

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"time"
)

type TaskComment struct {
	Id        string    `json:"id"`
	TaskId    string    `json:"task_id"`
	Message   string    `json:"message"`
	Author    int       `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type WsEvent struct {
	Type  string      `json:"type"`
	Event interface{} `json:"event"`
}

func LeaveComment(comment TaskComment) (sql.Result, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	comment.CreatedAt = time.Now()
	comment.Id = id.String()

	r, err := taskStorage.db.Exec("INSERT INTO task_comments (id, author, message, created_at, task_id) values (?,?,?,?,?)", comment.Id, comment.Author, comment.Message, comment.CreatedAt, comment.TaskId)

	if err != nil {
		log.Printf("Error while inserting task comment %v", err)
		return nil, err
	}

	// TODO add handle events by pipe not right here
	commentJson, err := json.Marshal(comment)
	_, err = taskStorage.db.Exec(`INSERT into tasks_events (task_id, event_type, payload, occurred_on) values (?, ?, ?, ?)`, comment.TaskId, "task_comment_left", commentJson, time.Now())

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}

	event := WsEvent{
		Type:  "comment_added",
		Event: comment,
	}
	wsEventJson, _ := json.Marshal(event)
	hub.Broadcast <- wsEventJson

	return r, nil
}

func LoadComments(taskId string, limit int) []TaskComment {
	taskCommentsList := make([]TaskComment, 0)

	rows, err := taskStorage.db.Query("SELECT id, task_id, message, author, created_at FROM task_comments WHERE task_id= ? LIMIT ?", taskId, limit)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var taskComment TaskComment
		err = rows.Scan(&taskComment.Id, &taskComment.TaskId, &taskComment.Message, &taskComment.Author, &taskComment.CreatedAt)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			taskCommentsList = append(taskCommentsList, taskComment)
		}
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return taskCommentsList
}
