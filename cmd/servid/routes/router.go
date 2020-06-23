package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"tasks17-server/cmd/servid/routes/handlers"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/storage"
)

func handle404() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaUrl, _ := url.Parse(r.Header.Get("Referer"))

		w.Header().Set("Access-Control-Allow-Origin", "http://"+spaUrl.Host)
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")

		if r.Method != "OPTIONS" {
			log.Println("[404] Try to perform url: ", r.URL)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	})
}

func handleMethodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaUrl, _ := url.Parse(r.Header.Get("Referer"))

		w.Header().Set("Access-Control-Allow-Origin", "http://"+spaUrl.Host)
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")

		if r.Method != "OPTIONS" {
			log.Println("[404] Try to perform url: ", r.URL)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}

func CreateRouter(h *platform.Hub, db *sql.DB) http.Handler {
	s := storage.New(db)
	authenticator := auth.NewAuthenticator(db)
	r := mux.NewRouter()

	r.Use(logMiddleware)
	r.Use(corsMiddleware)
	r.Use(authorizeMiddleware(&s))

	r.NotFoundHandler = handle404()
	r.MethodNotAllowedHandler = handleMethodNotAllowed()

	r.HandleFunc("/sign-in", handlers.SignInHandler(authenticator)).Methods("POST", "OPTIONS")

	r.HandleFunc("/projects", handlers.ProjectListHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/projects/add", handlers.AddProjectHandler(&s)).Methods("POST", "OPTIONS")
	r.HandleFunc("/project/{project}", handlers.ProjectInfoHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/project/{project}/stages/add", handlers.ProjectAddStageHandler(&s)).Methods("POST", "OPTIONS")
	r.HandleFunc("/project/{project}/stage/{stage}/tasks/add", handlers.AddTaskForProjectStageHandler(&s)).Methods("POST", "OPTIONS")

	r.HandleFunc("/teams", handlers.TeamListHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/teams/add", handlers.AddTeamHandler(&s)).Methods("POST", "OPTIONS")

	// Tasks package routes
	r.HandleFunc("/team/{team}", handlers.TeamInfoHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks", handlers.TaskListHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks/state", handlers.TaskStateListHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks/add", handlers.TaskCreateHandler(&s)).Methods("POST", "OPTIONS")

	r.HandleFunc("/team/{team}/tasks/{task}", handlers.TaskLoadHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks/{task}/close", handlers.TaskClose(&s)).Methods("POST", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks/{task}/update-description", handlers.TaskUpdateDescription(&s)).Methods("POST", "OPTIONS")

	r.HandleFunc("/team/{team}/tasks/{task}/comments/add", handlers.DiscussionCommentCreateHandler(h, &s)).Methods("POST", "OPTIONS")

	r.HandleFunc("/team/{team}/stages", handlers.StagesLoadHandler(&s)).Methods("GET", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks/{task}/add-sub-task", handlers.SubTaskCreateHandler(h, &s)).Methods("POST", "OPTIONS")
	r.HandleFunc("/team/{team}/tasks/{task}/{subTask}/close", handlers.SubTaskCloseHandler(&s)).Methods("POST", "OPTIONS")

	// Discussions
	r.HandleFunc("/discussions/{discussion}/comments/add", handlers.DiscussionCommentCreateHandler(h, &s)).Methods("POST", "OPTIONS")
	r.HandleFunc("/discussions/{discussion}/update-last-comment", handlers.UpdateLastCommentHandler(&s)).Methods("POST", "OPTIONS")

	// WebSocket Server
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		platform.ServeWs(w, r)
	})

	return r
}
