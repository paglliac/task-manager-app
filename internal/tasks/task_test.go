package tasks_test

import (
	_ "github.com/lib/pq"
	"tasks17-server/internal/tasks"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	// Task creation
	expectedTask := tasks.Task{
		Title:       "New Task",
		Description: "Description",
		Status:      "open",
		AuthorId:    setup.users.main().id,
		TeamId:      setup.users.main().teamId,
	}
	id, _ := tasks.CreateTask(&s, &expectedTask)
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

	states := tasks.LoadTaskStateList(&s, setup.users.secondary().id, setup.users.secondary().teamId)

	if len(states) != 1 {
		t.Errorf("Expected states 1, but got %d", len(states))
	}

	if states[id] == nil {
		t.Errorf("Expected states contains task state with task id %s", id)
	}
}

func TestReadComment(t *testing.T) {
	expectedTask := tasks.Task{
		Title:       "New Task",
		Description: "Description",
		Status:      "open",
		AuthorId:    setup.users.main().id,
		TeamId:      setup.users.main().teamId,
	}
	id, _ := tasks.CreateTask(&s, &expectedTask)
	defer s.RemoveTask(id)

	commentId, _ := tasks.LeaveComment(hh, &s, tasks.Comment{
		DiscussionId: expectedTask.DiscussionId,
		Message:      "New comment",
		CreatedAt:    time.Now(),
		Author:       setup.users.main().id,
	})

	states := tasks.LoadTaskStateList(&s, setup.users.secondary().id, setup.users.secondary().teamId)

	if states[id].UnreadComments != 1 {
		t.Errorf("Expected one unread comments, got %d", states[id].UnreadComments)
	}

	s.UpdateLastWatchedComment(setup.users.secondary().id, expectedTask.DiscussionId, commentId)

	states = tasks.LoadTaskStateList(&s, setup.users.secondary().id, setup.users.secondary().teamId)

	if states[id].UnreadComments != 0 {
		t.Errorf("Expected zero unread comments, got %d", states[id].UnreadComments)
	}
}
