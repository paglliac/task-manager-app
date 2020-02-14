package main

import (
	"ResearchGolang/cmd/servid/routes"
	"ResearchGolang/internal/models"
	"ResearchGolang/internal/platform"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

func main() {
	err := platform.InitDB()

	if err != nil {
		panic(err)
	}

	platform.InitHub()

	r := routes.CreateRouter()

	models.InitTasksModule(platform.Db)

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
