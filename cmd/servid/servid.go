package main

import (
	"ResearchGolang/cmd/servid/routes/handlers"
	"ResearchGolang/internal/platform"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func LoggedHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		log.Println(r.URL)
		handler(w, r)
	})
}

func main() {
	err := platform.InitDB()

	if err != nil {
		panic(err)
	}

	hub := platform.InitHub()
	go hub.Run()

	LoggedHandleFunc("/", handlers.UsersListHandler)

	LoggedHandleFunc("/messages", handlers.MessageListHandler)
	LoggedHandleFunc("/messages/add", handlers.MessageCreateHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		platform.ServeWs(hub, w, r)
	})

	log.Println("Server have been started listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
