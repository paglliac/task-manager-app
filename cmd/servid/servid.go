package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"tasks17-server/cmd/servid/routes"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/tasks"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	config := platform.InitConfig()
	db, err := platform.InitDb(config.DbConfig)
	if err != nil {
		panic(err)
	}

	hub := platform.InitHub()

	tasks.Init(hub)
	auth.Init(db)

	r := routes.CreateRouter(hub, db)

	log.Println("Server have been started listening on port 8080")

	srv := http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	log.Fatal(srv.ListenAndServe())
}
