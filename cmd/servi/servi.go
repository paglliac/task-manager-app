package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"tasks17-server/internal/platform"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	db, err := platform.InitDB(dbHost)

	if err != nil {
		panic(err)
	}

	cmd := os.Args[1]
	log.Println(cmd)

	if cmd == "createUser" {
		createUserCmd := flag.NewFlagSet("createUser", flag.ExitOnError)
		username := createUserCmd.String("username", "example", "a string")
		email := createUserCmd.String("email", "", "a string")
		password := createUserCmd.String("password", "", "a string")

		createUserCmd.Parse(os.Args[2:])

		var uid int

		err := db.QueryRow(`INSERT into users (name, email, password) values ($1, $2, $3) RETURNING id`, *username, *email, *password).Scan(&uid)

		if err != nil {
			log.Printf("[save user] error while saving user: %v", err)
			os.Exit(1)
		}

		log.Printf("User created successfully with id: %d", uid)
		return
	}

	fmt.Printf("Unknown command %s \n", cmd)
	os.Exit(127)
}
