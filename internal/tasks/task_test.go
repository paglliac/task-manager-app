package tasks

import (
	"ResearchGolang/internal/platform"
	_ "github.com/go-sql-driver/mysql"
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
