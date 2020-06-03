package tasks

import (
	"tasks17-server/internal/platform"
	"time"
)

type SubTask struct {
	Id        int       `json:"id"`
	TaskId    string    `json:"task_id"`
	StageId   int       `json:"stage_id"`
	Rank      int       `json:"rank"`
	AuthorId  int       `json:"author_id"`
	Status    int       `json:"status"`
	Name      string    `json:"name"`
	ClosedAt  time.Time `json:"closed_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Stage struct {
	Id     int
	TeamId int
	Name   string
	Rank   int
}

type Progress struct {
	SubTasks []SubTask `json:"sub_tasks"`
	Stages   []Stage   `json:"stages"`
}

func AddSubTask(h *platform.Hub, ts TaskStorage, subTask SubTask) (int, error) {
	subTask.CreatedAt = time.Now()

	id, err := ts.SaveSubTask(subTask)

	h.Handle(platform.WsEvent{Type: "sub_task_added", Event: subTask})

	return id, err
}

func CompleteSubTask(ts TaskStorage, id int) {
	ts.CompleteSubTask(id)

	hub.Handle(platform.WsEvent{Type: "sub_task_completed", Event: id})
}

func LoadProgress(ts TaskStorage, taskId string) Progress {
	stList := ts.LoadSubTasks(taskId)
	stageList := ts.LoadTaskStages(taskId)

	return Progress{
		SubTasks: stList,
		Stages:   stageList,
	}
}
