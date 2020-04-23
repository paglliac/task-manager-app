package tasks

import (
	"tasks17-server/internal/platform"
	"time"
)

type TaskComment struct {
	Id         string    `json:"id"`
	TaskId     string    `json:"task_id"`
	Message    string    `json:"message"`
	Author     int       `json:"author"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}

func LeaveComment(h *platform.Hub, ts TaskStorage, comment TaskComment) (id string, err error) {
	id, err = ts.SaveComment(comment)

	if err != nil {
		return "", err
	}

	ts.SaveCommentEvent(comment)
	ts.UpdateLastWatchedComment(comment.Author, comment.TaskId, comment.Id)

	h.Handle(platform.WsEvent{Type: "comment_added", Event: comment})

	return id, nil
}
