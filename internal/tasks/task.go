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
	taskId     int
	eventType  string
	occurredOn string
}

type TaskState struct {
	taskId         string
	taskName       string
	unreadComments string
}

type TaskStorage interface {
	loadTasks(limit int) []Task
	saveTask(task Task) (sql.Result, error)
	loadTasksState(userId int) []TaskState
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

func LoadTaskStateList(userId int) []TaskState {
	return taskStorage.loadTasksState(userId)
}

func LoadTasks(limit int) []Task {
	return taskStorage.loadTasks(limit)
}

func CreateTask(task Task) (sql.Result, error) {
	return taskStorage.saveTask(task)
}
