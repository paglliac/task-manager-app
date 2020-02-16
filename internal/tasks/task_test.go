package tasks

import (
	_ "github.com/go-sql-driver/mysql"
	"tasks17-server/internal/platform"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	platform.InitDB()
	_, err := CreateTask(Task{
		Title:       "New Task",
		Description: "Description",
		Status:      taskStatusOpen,
		CreatedAt:   time.Time{},
	})

	if err != nil {
		t.Error(err)
	}
}
