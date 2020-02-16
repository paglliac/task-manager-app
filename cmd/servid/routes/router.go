package routes

import (
	"tasks17-server/cmd/servid/routes/handlers"
	"tasks17-server/internal/platform"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

		next.ServeHTTP(w, r)
	})
}

func CreateRouter() http.Handler {
	r := mux.NewRouter()

	r.Use(CorsMiddleware)
	r.Use(LogMiddleware)

	r.HandleFunc("/", handlers.UsersListHandler).Methods("GET", "OPTIONS")

	// Tasks package routes
	r.HandleFunc("/tasks", handlers.TaskListHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/add", handlers.TaskCreateHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks/{comment}/comments", handlers.TaskCommentLoadHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/comments/add", handlers.TaskCommentCreateHandler).Methods("POST", "OPTIONS")

	// WebSocket Server
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		platform.ServeWs(w, r)
	})

	return r
}
