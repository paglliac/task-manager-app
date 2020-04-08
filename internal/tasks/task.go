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

const (
	taskStatusOpen = "open"
)

func (t Task) UpdateDescription(description string) error {
	return taskStorage.updateDescription(t.Id, description)
}

func (t Task) Close() error {
	return taskStorage.closeTask(t.Id)
}

func (t Task) CompleteSubTask(subTask int) {
	_, _ = taskStorage.getDb().Exec("UPDATE sub_tasks SET status = $1 WHERE id = $2", 1, subTask)

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
	loadTask(id string) Task
	saveTask(t Task) (sql.Result, error)
	saveSubTask(st SubTask) (sql.Result, error)
	loadStates(uid int) map[string]*State
	loadEvents(uid int) []Event
	// TODO hack for not refactoring task_comment db interactions need fix asap
	getDb() *platform.Storage
	updateDescription(id string, description string) error
	closeTask(id string) error
}

var taskStorage TaskStorage

var hub *platform.Hub

func InitTasksModule(db platform.Storage, h *platform.Hub) {
	taskStorage = &SqlTaskStorage{db: db}
	hub = h
}

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

func AddSubTask(subTask SubTask) (sql.Result, error) {
	subTask.CreatedAt = time.Now()
	return taskStorage.saveSubTask(subTask)
}
