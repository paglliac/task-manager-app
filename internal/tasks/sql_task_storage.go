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

func (s *SqlTaskStorage) loadStates(userId int) map[string]*TaskState {
	rows, err := s.db.Query("SELECT id, title FROM tasks")
	states := make(map[string]*TaskState)

	// This need to avoid unused parameter, but user id not used at the moment
	log.Println(userId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var s TaskState
		err = rows.Scan(&s.TaskId, &s.TaskTitle)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			states[s.TaskId] = &s
		}
	}

	return states
}

func (s *SqlTaskStorage) getDb() *platform.Storage {
	return &s.db
}

func (s *SqlTaskStorage) loadEvents(userId int) []TaskEvent {
	events := make([]TaskEvent, 0)
	rows, err := s.db.Query(`SELECT te.task_id, te.event_type, te.payload, te.occurred_on 
										FROM tasks_events te
										LEFT JOIN task_last_watched_event tlwe ON te.task_id = tlwe.task_id AND tlwe.user_id = ?
									WHERE te.id > tlwe.last_event_id or tlwe.last_event_id IS NULL`, userId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var payload sql.NullString
		var e TaskEvent

		err = rows.Scan(&e.taskId, &e.eventType, &payload, &e.occurredOn)
		e.payload = payload.String

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			events = append(events, e)
		}
	}

	return events
}

func (s *SqlTaskStorage) loadTask(taskId string) Task {
	row := s.db.QueryRow("SELECT id, title, description, status, created_at FROM tasks WHERE id = ?", taskId)

	var task Task

	err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)

	if err != nil {
		log.Println("Error while scanning entity", err)
	}

	return task
}

func (s *SqlTaskStorage) loadTasks(limit int) []Task {
	taskList := make([]Task, 0)

	rows, err := s.db.Query("SELECT id, title, description, status, created_at FROM tasks limit ?", limit)
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

	r, err := s.db.Exec(`INSERT into tasks (id, title, description, status, created_at, updated_at, author) values (?, ?, ?, ?, ?, ?, ?)`, id, task.Title, task.Description, taskStatusOpen, time.Now(), time.Now(), task.AuthorId)

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
