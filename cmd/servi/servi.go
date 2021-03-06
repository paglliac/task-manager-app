package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/storage"
)

func main() {
	config := platform.InitConfig()
	db, err := platform.InitDb(config.DbConfig)
	s := storage.New(db)

	if err != nil {
		panic(err)
	}

	cmd := os.Args[1]
	log.Println(cmd)

	if cmd == "createUser" {
		createUser(s)
		return
	}

	if cmd == "setup" {
		setup(s)
		return
	}

	fmt.Printf("Unknown command %s \n", cmd)
	os.Exit(127)
}

func createUser(db storage.Storage) {
	createUserCmd := flag.NewFlagSet("createUser", flag.ExitOnError)
	orgId := createUserCmd.String("org", "1", "an int")
	username := createUserCmd.String("username", "example", "a string")
	email := createUserCmd.String("email", "", "a string")
	password := createUserCmd.String("password", "", "a string")

	err := createUserCmd.Parse(os.Args[2:])
	if err != nil {
		log.Printf("[save user] parsing cmd parameters error: %v", err)
		os.Exit(1)
	}

	var uid int

	err = db.QueryRow(`INSERT into users (organisation_id, name, email, password) values ($1, $2, $3, $4) RETURNING id`, *orgId, *username, *email, *password).Scan(&uid)

	if err != nil {
		log.Printf("[save user] error while saving: %v", err)
		os.Exit(1)
	}

	log.Printf("User created successfully with id: %d", uid)
}

func setup(db storage.Storage) {
	var orgId int
	var teamId int
	var uId int

	err := db.QueryRow(`INSERT into organisations (name) values ('Gangster Elephant') RETURNING id`).Scan(&orgId)

	if err != nil {
		log.Printf("[organisation saving] error while saving: %v", err)
		os.Exit(1)
	}

	log.Printf("Organisation created successfully with id: %d", orgId)

	err = db.QueryRow(`INSERT into teams (name, organisation_id) values ('First team', $1) RETURNING id`, orgId).Scan(&teamId)

	if err != nil {
		log.Printf("[team saving] error while saving: %v", err)
		os.Exit(1)
	}

	log.Printf("Team created successfully with id: %d", teamId)

	err = db.QueryRow(`INSERT into sub_task_stages (team_id,rank, name) values ($1, 1, 'Development') RETURNING id`, teamId).Scan(&uId)
	log.Printf("SubTaskStage Development created successfully with id: %d", uId)

	err = db.QueryRow(`INSERT into sub_task_stages (team_id,rank, name) values ($1, 2, 'QA') RETURNING id`, teamId).Scan(&uId)
	log.Printf("SubTaskStage QA created successfully with id: %d", uId)

	err = db.QueryRow(`INSERT into users (organisation_id, name, email, password) values ($1, 'kir', 'kir@gangsterelephant.io', 'kir') RETURNING id`, orgId).Scan(&uId)
	log.Printf("User created successfully with id: %d", uId)

	err = db.QueryRow(`INSERT into users (organisation_id, name, email, password) values ($1, 'serg', 'serg@gangsterelephant.io', 'serg') RETURNING id`, orgId).Scan(&uId)
	log.Printf("User created successfully with id: %d", uId)
}
