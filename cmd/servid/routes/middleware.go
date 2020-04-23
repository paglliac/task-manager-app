package routes

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/tasks"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO need to change spa url with env variable
		spaUrl, _ := url.Parse(r.Header.Get("Referer"))
		w.Header().Set("Access-Control-Allow-Origin", "http://"+spaUrl.Host)
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authorizeMiddleware(ts tasks.TaskStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := r.Header.Get("Authorization")
			credentials := auth.ParseToken(t)

			teamId, err := strconv.Atoi(mux.Vars(r)["team"])

			if err == nil {
				team := ts.LoadTeam(teamId)

				if team.OrgId != credentials.Oid {
					http.Error(w, "Wrong team id", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(auth.NewContextWithCredentials(r.Context(), credentials)))
		})
	}
}
