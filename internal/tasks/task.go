package tasks

import (
	"database/sql"
	"tasks17-server/internal/platform"
	"time"
)

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskEvent struct {
	taskId     string
	eventType  string
	occurredOn time.Time
	payload    string
}

type TaskState struct {
	TaskId         string `json:"task_id"`
	TaskName       string `json:"task_name"`
	UnreadComments int    `json:"unread_comments"`
}

type TaskStorage interface {
	loadTasks(limit int) []Task
	saveTask(task Task) (sql.Result, error)
	loadTasksState(userId int) (map[string]*TaskState, []TaskEvent)
	// TODO hack for not refactoring task_comment db interactions need fix asap
	getDb() *platform.Storage
}

var taskStorage TaskStorage

var hub *platform.Hub

func InitTasksModule(db platform.Storage, h *platform.Hub) {
	taskStorage = &SqlTaskStorage{db: db}
	hub = h
}

const (
	taskStatusOpen = "open"
)

func LoadTaskStateList(userId int) map[string]*TaskState {
	states, events := taskStorage.loadTasksState(userId)
	for _, event := range events {
		if event.eventType == "task_comment_left" {
			states[event.taskId].UnreadComments++
		}
	}

	return states
}

func LoadTasks(limit int) []Task {
	return taskStorage.loadTasks(limit)
}

func CreateTask(task Task) (sql.Result, error) {
	return taskStorage.saveTask(task)
}
