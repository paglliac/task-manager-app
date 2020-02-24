package tasks

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"time"
)

type TaskComment struct {
	Id         string    `json:"id"`
	TaskId     string    `json:"task_id"`
	Message    string    `json:"message"`
	Author     int       `json:"author"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
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

	r, err := taskStorage.getDb().Exec("INSERT INTO task_comments (id, author, message, created_at, task_id) values (?,?,?,?,?)", comment.Id, comment.Author, comment.Message, comment.CreatedAt, comment.TaskId)

	if err != nil {
		log.Printf("Error while inserting task comment %v", err)
		return nil, err
	}

	// TODO add handle events by pipe not right here
	commentJson, err := json.Marshal(comment)
	_, err = taskStorage.getDb().Exec(`INSERT into tasks_events (task_id, event_type, payload, occurred_on) values (?, ?, ?, ?)`, comment.TaskId, "task_comment_left", commentJson, time.Now())

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

func LoadComments(taskId string) []TaskComment {
	taskCommentsList := make([]TaskComment, 0)

	rows, err := taskStorage.getDb().Query("SELECT task_comments.id, task_id, message, author, users.name, created_at FROM task_comments LEFT JOIN users ON users.id = task_comments.author WHERE task_id= ? ORDER BY created_at", taskId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var taskComment TaskComment
		err = rows.Scan(&taskComment.Id, &taskComment.TaskId, &taskComment.Message, &taskComment.Author, &taskComment.AuthorName, &taskComment.CreatedAt)

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

func UpdateLastWatchedComment(userId int, taskId string, commentId string) {
	id := findLastEventByCommentId(commentId)
	log.Println(id)
	_, err := taskStorage.getDb().Exec(`INSERT into task_last_watched_event (user_id, task_id, last_event_id) values (?, ?, ?) ON DUPLICATE KEY UPDATE last_event_id = ?`, userId, taskId, id, id)

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}
}

func findLastEventByCommentId(id string) int {
	rows, err := taskStorage.getDb().Query("SELECT id FROM tasks_events WHERE event_type = 'task_comment_left' AND payload LIKE ? LIMIT 1", "%"+id+"%")
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var eventId int
		err = rows.Scan(&eventId)

		if err != nil {
			log.Println("Error while scanning entity", err)
		}
		return eventId
	}

	return 0
}
