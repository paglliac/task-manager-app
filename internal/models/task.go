package models

import (
	"ResearchGolang/internal/platform"
	"fmt"
	"log"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func LoadTasks(limit int) []Task {
	taskList := make([]Task, 0)

	q := fmt.Sprintf("SELECT id, title, description, status, created_at FROM tasks limit %d", limit)
	rows, err := platform.Db.Query(q)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			taskList = append(taskList, task)
		}
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return taskList
}
