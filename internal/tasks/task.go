package tasks

import (
	"database/sql"
	"log"
	"tasks17-server/internal/platform"
	"time"
)

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorId    int
}

func (t Task) UpdateDescription(description string) error {
	return taskStorage.updateDescription(t.Id, description)
}

func (t Task) Close() error {
	return taskStorage.closeTask(t.Id)
}

type Event struct {
	taskId     string
	eventType  string
	occurredOn time.Time
	payload    string
}

type State struct {
	TaskId         string `json:"task_id"`
	TaskTitle      string `json:"task_title"`
	UnreadComments int    `json:"unread_comments"`
}

type TaskStorage interface {
	loadTasks(limit int) []Task
	loadTask(taskId string) Task
	saveTask(task Task) (sql.Result, error)
	loadStates(userId int) map[string]*State
	loadEvents(userId int) []Event
	// TODO hack for not refactoring task_comment db interactions need fix asap
	getDb() *platform.Storage
	updateDescription(taskId string, description string) error
	closeTask(id string) error
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

func LoadTaskStateList(userId int) map[string]*State {
	events := taskStorage.loadEvents(userId)
	states := taskStorage.loadStates(userId)

	for _, event := range events {
		if event.eventType == "task_comment_left" {
			if _, ok := states[event.taskId]; ok {
				states[event.taskId].UnreadComments++
			} else {
				log.Println("ERR event for not exists task")
			}
		}
	}

	return states
}

func LoadTask(taskId string) Task {
	return taskStorage.loadTask(taskId)
}

func LoadTasks(limit int) []Task {
	return taskStorage.loadTasks(limit)
}

func CreateTask(task Task) (sql.Result, error) {
	return taskStorage.saveTask(task)
}
