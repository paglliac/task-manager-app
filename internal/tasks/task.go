package tasks

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
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
}

type TaskEvent struct {
	taskId     int
	eventType  string
	occurredOn string
}

type TaskUpdate struct {
	taskId         string
	taskName       string
	unreadComments string
}

type TaskStorage interface {
	loadTasks() []Task
}

var taskStorage *SqlTaskStorage

var hub *platform.Hub

func InitTasksModule(db platform.Storage, h *platform.Hub) {
	taskStorage = &SqlTaskStorage{db: db}
	hub = h
}

type SqlTaskStorage struct {
	db platform.Storage
}

func (s *SqlTaskStorage) loadTasks(limit int) []Task {
	taskList := make([]Task, 0)

	q := fmt.Sprintf("SELECT id, title, description, status, created_at FROM tasks limit %d", limit)
	rows, err := s.db.Query(q)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			taskList = append(taskList, task)
		}
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return taskList
}

const (
	taskStatusOpen = "open"
)

func LoadTasks(limit int) []Task {
	return taskStorage.loadTasks(limit)
}

func CreateTask(task Task) (sql.Result, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Printf("Error while uuid generated. %v \n", err)
		return nil, fmt.Errorf("error while generated uuid %v", err)
	}

	r, err := taskStorage.db.Exec(`INSERT into tasks (id, title, description, status, created_at, updated_at, author) values (?, ?, ?, ?, ?, ?, ?)`, id, task.Title, task.Description, taskStatusOpen, time.Now(), time.Now(), 1)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// TODO add handle events by pipe not right here
	_, err = taskStorage.db.Exec(`INSERT into tasks_events (task_id, event_type, occurred_on) values (?, ?, ?)`, id, "task_created", time.Now())

	if err != nil {
		log.Println(err)
	}

	return r, nil
}
