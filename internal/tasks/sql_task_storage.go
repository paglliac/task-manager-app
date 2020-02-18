package tasks

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"tasks17-server/internal/platform"
	"time"
)

type SqlTaskStorage struct {
	db platform.Storage
}

func (s *SqlTaskStorage) loadTasks(limit int) []Task {
	taskList := make([]Task, 0)

	q := fmt.Sprintf("SELECT id, title, description, status, created_at FROM tasks limit %d", limit)
	rows, err := s.db.Query(q)
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

func (s *SqlTaskStorage) saveTask(task Task) (sql.Result, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Printf("Error while uuid generated. %v \n", err)
		return nil, fmt.Errorf("error while generated uuid %v", err)
	}

	r, err := s.db.Exec(`INSERT into tasks (id, title, description, status, created_at, updated_at, author) values (?, ?, ?, ?, ?, ?, ?)`, id, task.Title, task.Description, taskStatusOpen, time.Now(), time.Now(), 1)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// TODO add handle events by pipe not right here
	_, err = s.db.Exec(`INSERT into tasks_events (task_id, event_type, occurred_on) values (?, ?, ?)`, id, "task_created", time.Now())

	if err != nil {
		log.Println(err)
	}

	return r, nil
}
