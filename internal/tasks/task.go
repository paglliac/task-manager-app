package tasks

import (
	"github.com/google/uuid"
	"tasks17-server/internal/platform"
	"time"
)

var hub EventHandler

func Init(h EventHandler) {
	hub = h
}

type Task struct {
	Id          string    `json:"id"`
	TeamId      int       `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AuthorId    int
}

func (t *Task) PreSave() {
	if t.Id == "" {
		t.Id = uuid.New().String()
	}

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}

	t.UpdatedAt = time.Now()
}

func CreateTask(ts TaskStorage, task Task) (string, error) {
	err := ts.SaveTask(&task)
	hub.Handle(platform.WsEvent{Type: "task_added", Event: task.Id})
	return task.Id, err
}
