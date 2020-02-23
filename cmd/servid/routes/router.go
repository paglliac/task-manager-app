package routes

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"tasks17-server/cmd/servid/routes/handlers"
	"tasks17-server/internal/platform"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, _ := url.Parse(r.Header.Get("Referer"))
		w.Header().Set("Access-Control-Allow-Origin", "http://"+url.Host)
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CreateRouter() http.Handler {
	r := mux.NewRouter()

	r.Use(LogMiddleware)
	r.Use(CorsMiddleware)

	r.HandleFunc("/sign-in", handlers.SignInHandler).Methods("POST", "OPTIONS")

	// Tasks package routes
	r.HandleFunc("/tasks", handlers.TaskListHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/state", handlers.TaskStateListHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/add", handlers.TaskCreateHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks/update-last-event", handlers.TaskUpdateLastCommentHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks/{task}", handlers.TaskLoadHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/comments/add", handlers.TaskCommentCreateHandler).Methods("POST", "OPTIONS")

	// WebSocket Server
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		platform.ServeWs(w, r)
	})

	return r
}
