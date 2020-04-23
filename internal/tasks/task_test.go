package tasks_test

import (
	_ "github.com/lib/pq"
	"tasks17-server/internal/tasks"
	"testing"
)

func TestCreateTask(t *testing.T) {
	// Task creation
	expectedTask := tasks.Task{
		Title:       "New Task",
		Description: "Description",
		Status:      "open",
		AuthorId:    1,
		TeamId:      1,
	}
	id, _ := tasks.CreateTask(&s, expectedTask)
	defer s.RemoveTask(id)

	// Check task has been created
	task := s.LoadTask(id)

	expectedTask.CreatedAt = task.CreatedAt
	expectedTask.UpdatedAt = task.UpdatedAt
	expectedTask.Id = task.Id

	if task != expectedTask {
		t.Errorf("Task not created correctly \n %+v \n %+v", expectedTask, task)
	}

	// Check event has been thrown
	if len(h.events) != 1 {
		t.Error("WebSocket event not been thrown")
	}

	event := h.events[0]
	if event.Type != "task_added" || event.Event != task.Id {
		t.Error("Incorrect event has been thrown")
	}
}
