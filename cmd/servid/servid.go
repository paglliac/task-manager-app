package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"tasks17-server/cmd/servid/routes"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/tasks"
	"time"
)

func main() {
	db, err := platform.InitDB()

	if err != nil {
		panic(err)
	}

	hub := platform.InitHub()

	r := routes.CreateRouter()

	tasks.InitTasksModule(db, hub)
	auth.InitAuthModule(db)

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
