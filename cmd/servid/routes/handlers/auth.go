package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tasks17-server/internal/auth"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInHandler(auth auth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ar AuthRequest

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ar)

		if err != nil {
			http.Error(w, "Error while decode body", http.StatusUnprocessableEntity)
			return
		}

		jwt, err := auth.IssueToken(ar.Email, ar.Password)

		if err != nil {
			http.Error(w, "Wrong credentials", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "{\"token\": \"%s\"}", jwt)
	}
}
