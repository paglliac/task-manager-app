package handlers

import (
	"net/http"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/mail"
	"tasks17-server/internal/tasks"
)

func UserProfileHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		var rs response

		credentials, _ := auth.FromRequest(r)
		u := ts.LoadUser(credentials.Uid)

		mail.Send(mail.DefaultMail{
			Subject: "How are you?",
			Body:    "How are you my old best friend",
			To:      "kir@alyce.com",
		})

		rs.Name = u.Name
		rs.Email = u.Email

		jsonResponse(rs, w)
	}
}
