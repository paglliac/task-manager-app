package routes

import (
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"tasks17-server/cmd/servid/routes/handlers"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/tasks"
	"testing"
)

func TestSignIn(t *testing.T) {
	ar := handlers.AuthRequest{
		Email:    setup.user.email,
		Password: setup.user.password,
	}

	expectedCredentials := auth.Credentials{
		Uid: setup.user.id,
		Oid: setup.orgId,
	}

	response := NewRequester(t).post("/sign-in", ar)

	c := auth.ParseToken(response.get("token"))

	if *c != expectedCredentials {
		t.Errorf("Credentials is invalid")
	}
}

func TestTaskListEndpoint(t *testing.T) {
	setup.user.createTask()
	setup.user.createTask()

	response := NewRequester(t).auth(&setup.user).get(fmt.Sprintf("/team/%d/tasks", setup.teamId))

	tt := response.getRaw("tasks")

	if len(tt.([]interface{})) != 2 {
		t.Errorf("Not all tasks loaded. Expected %d, but got %d", 2, len(tt.([]interface{})))
	}
}

func TestTasksStateListEndpoint(t *testing.T) {
	setup.user.createTask()
	setup.user.createTask()

	// TODO add testing for unread comments
	tasks := s.LoadTasks(setup.teamId)
	response := NewRequester(t).auth(&setup.user).get(fmt.Sprintf("/team/%d/tasks/state", setup.teamId))

	tt := response.getRaw("tasks")

	if len(tt.(map[string]interface{})) != len(tasks) {
		t.Errorf("Not all tasks loaded. Expected %d, got %d", len(tasks), len(tt.(map[string]interface{})))
	}
}

func TestTeamsEndpoint(t *testing.T) {
	response := NewRequester(t).auth(&setup.user).get("/teams")

	bytes, _ := ioutil.ReadAll(response.rawResponse.Body)
	r := string(bytes)

	if r != fmt.Sprintf("[{\"id\":%d,\"org_id\":%d,\"name\":\"%s\"}]", setup.teamId, setup.orgId, setup.teamName) {
		t.Error("Response incorrect")
	}
}

func TestStagesEndpoint(t *testing.T) {
	response := NewRequester(t).auth(&setup.user).get(fmt.Sprintf("/team/%d/stages", setup.teamId))

	bytes, _ := ioutil.ReadAll(response.rawResponse.Body)
	r := string(bytes)

	if r != "{\"stages\":[]}" {
		t.Errorf("Response incorrect. Expected %s, got %s", "{\"stages\":[]}", r)
	}
}

func TestTaskCreateEndpoint(t *testing.T) {
	tr := map[string]interface{}{
		"title":       "New task",
		"description": "New task description",
	}

	response := NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/add", setup.teamId), tr)

	taskId := response.get("id")

	if taskId == "0" || taskId == "" {
		t.Error("Task have not bee created or id not generated")
		t.FailNow()
	}

	s.RemoveTask(taskId)
}

func TestTaskLoadEndpoint(t *testing.T) {
	task, rollback := setup.user.createTask()
	defer rollback()

	response := NewRequester(t).auth(&setup.user).get(fmt.Sprintf("/team/%d/tasks/%s", setup.teamId, task.Id))

	decodedTask := response.getRaw("task").(map[string]interface{})

	if task.Id != decodedTask["id"] {
		t.Error("Task not loaded")
		t.FailNow()
	}
}

func TestTaskCloseEndpoint(t *testing.T) {
	// TODO create new task every time
	task, rollback := setup.user.createTask()
	defer rollback()

	NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/close", setup.teamId, task.Id), nil)

	updatedTask := s.LoadTask(task.Id)

	if updatedTask.Status != "closed" {
		t.Error("Task have not been closed")
	}
}

func TestTaskChangeDescriptionEndpoint(t *testing.T) {
	// TODO check restriction to update closed task
	task, rollback := setup.user.createTask()
	defer rollback()

	description := map[string]string{
		"description": "new_description",
	}

	NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/update-description", setup.teamId, task.Id), description)

	updatedTask := s.LoadTask(task.Id)

	if updatedTask.Description != description["description"] {
		t.Error("Task have not been closed")
	}
}

func TestAddCommentEndpoint(t *testing.T) {
	task, rollback := setup.user.createTask()
	defer rollback()

	response := NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/comments/add", setup.teamId, task.Id), tasks.TaskComment{
		TaskId:  task.Id,
		Message: "New comment",
		Author:  setup.user.id,
	})

	commentId := response.get("id")

	if commentId == "" {
		t.Error("Comment id not returned")
	}

	comments := s.LoadComments(task.Id)

	if len(comments) == 0 {
		t.Error("Comment not created")
	}
}

func TestUpdateLastWatchedComment(t *testing.T) {
	task, rollback := setup.user.createTask()
	defer rollback()

	response := NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/comments/add", setup.teamId, task.Id), tasks.TaskComment{
		TaskId:  task.Id,
		Message: "New comment",
		Author:  setup.user.id,
	})

	commentId := response.get("id")

	NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/update-last-event", setup.teamId, task.Id), map[string]string{
		"comment_id": commentId,
	})
	// TODO add logic for unread comments
}

func TestAddSubTask(t *testing.T) {
	task, rollback := setup.user.createTask()
	defer rollback()

	stageId, _ := s.SaveStage(tasks.Stage{
		TeamId: setup.teamId,
		Name:   "Development",
		Rank:   0,
	})

	response := NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/add-sub-task", setup.teamId, task.Id), tasks.SubTask{
		TaskId:   task.Id,
		StageId:  stageId,
		AuthorId: setup.user.id,
		Status:   0,
		Name:     "New sub task",
	})

	subTaskId := response.getRaw("id").(float64)

	if subTaskId == 0 {
		t.Error("Sub task id not returned")
	}
}

func TestCompleteSubTask(t *testing.T) {
	task, rollback := setup.user.createTask()
	defer rollback()

	stageId, _ := s.SaveStage(tasks.Stage{
		TeamId: setup.teamId,
		Name:   "Development",
		Rank:   0,
	})

	response := NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/add-sub-task", setup.teamId, task.Id), tasks.SubTask{
		TaskId:   task.Id,
		StageId:  stageId,
		AuthorId: setup.user.id,
		Status:   0,
		Name:     "New sub task",
	})

	subTaskId := response.getRaw("id").(float64)

	NewRequester(t).auth(&setup.user).post(fmt.Sprintf("/team/%d/tasks/%s/%d/close", setup.teamId, task.Id, int64(subTaskId)), nil)
	subTasks := s.LoadSubTasks(task.Id)
	if subTasks[0].ClosedAt.IsZero() || subTasks[0].Status != 1 {
		t.Error("Sub task have not been closed")
	}
}
