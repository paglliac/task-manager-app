package tasks

import (
	"tasks17-server/internal/platform"
	"time"
)

type Discussion struct {
	id int
}

type Comment struct {
	Id           int       `json:"id"`
	DiscussionId string    `json:"-"`
	Message      string    `json:"message"`
	Author       int       `json:"author"`
	AuthorName   string    `json:"author_name"`
	CreatedAt    time.Time `json:"created_at"`
}

func LeaveComment(h *platform.Hub, ts TaskStorage, comment Comment) (id int, err error) {
	id, err = ts.SaveComment(comment)
	comment.Id = id

	if err != nil {
		return 0, err
	}

	ts.UpdateLastWatchedComment(comment.Author, comment.DiscussionId, comment.Id)

	h.Handle(platform.WsEvent{Type: "comment_added", Event: comment})

	return id, nil
}
