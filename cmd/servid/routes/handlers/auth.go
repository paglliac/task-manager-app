package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tasks17-server/internal/auth"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	type authRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var ar authRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ar)

	if err != nil {
		log.Println("Unprocessable entity")
	}

	id, err := auth.CheckAuth(ar.Email, ar.Password)

	fmt.Fprintf(w, "{\"id\": %d}", id)
}
