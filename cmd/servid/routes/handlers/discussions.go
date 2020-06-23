package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/tasks"
	"time"
)

func DiscussionCommentCreateHandler(h *platform.Hub, ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		discussionId := mux.Vars(r)["discussion"]
		credentials, _ := auth.FromRequest(r)
		var comment tasks.Comment

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&comment)

		if err != nil {
			http.Error(w, "Bad body", http.StatusBadRequest)
		}

		comment.Author = credentials.Uid
		comment.DiscussionId = discussionId
		comment.CreatedAt = time.Now()

		id, err := tasks.LeaveComment(h, ts, comment)

		if err != nil {
			http.Error(w, "Something went wrong", 500)
			return
		}

		jsonResponse(map[string]int{"id": id}, w)
	}
}

func UpdateLastCommentHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type updateLastComment struct {
			CommentId int `json:"comment_id"`
		}

		discussionId := mux.Vars(r)["discussionId"]
		credentials, _ := auth.FromRequest(r)
		var t updateLastComment

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&t)

		if err != nil {
			log.Printf("Error while decoding request: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ts.UpdateLastWatchedComment(credentials.Uid, discussionId, t.CommentId)
	}
}
