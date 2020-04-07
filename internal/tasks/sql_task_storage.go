package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"tasks17-server/internal/platform"
	"time"
)

type SqlTaskStorage struct {
	db platform.Storage
}

func (s *SqlTaskStorage) closeTask(id string) error {
	r, err := s.db.Exec("UPDATE tasks SET status = $1 WHERE id = $2", "closed", id)

	if err != nil {
		return err
	}

	if affected, _ := r.RowsAffected(); affected != 1 {
		return errors.New(fmt.Sprintf("task with id: %s not found", id))
	}

	return nil
}

func (s *SqlTaskStorage) updateDescription(taskId string, description string) error {
	r, err := s.db.Exec("UPDATE tasks SET description = $1 WHERE id = $2", description, taskId)

	if err != nil {
		return err
	}

	if affected, _ := r.RowsAffected(); affected != 1 {
		return errors.New(fmt.Sprintf("task with id: %s not found", taskId))
	}

	return nil
}

func (s *SqlTaskStorage) loadStates(userId int) map[string]*State {
	rows, err := s.db.Query("SELECT id, title FROM tasks where status = 'open'")
	states := make(map[string]*State)

	// This need to avoid unused parameter, but user id not used at the moment
	log.Println(userId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var s State
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

func (s *SqlTaskStorage) loadEvents(userId int) []Event {
	events := make([]Event, 0)
	rows, err := s.db.Query(`SELECT te.task_id, te.event_type, te.payload, te.occurred_on 
										FROM tasks_events te
										LEFT JOIN task_last_watched_event tlwe ON te.task_id = tlwe.task_id AND tlwe.user_id = $1
									WHERE te.id > tlwe.last_event_id or tlwe.last_event_id IS NULL`, userId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var payload sql.NullString
		var e Event

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
	row := s.db.QueryRow("SELECT id, title, description, status, created_at FROM tasks WHERE id = $1", taskId)

	var task Task

	err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)

	if err != nil {
		log.Println("Error while scanning entity", err)
	}

	return task
}

func (s *SqlTaskStorage) loadTasks(limit int) []Task {
	taskList := make([]Task, 0)

	rows, err := s.db.Query("SELECT id, title, description, status, created_at FROM tasks limit $1", limit)
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

	r, err := s.db.Exec(`INSERT into tasks (id, title, description, status, created_at, updated_at, author) values ($1, $2, $3, $4, $5, $6, $7)`, id, task.Title, task.Description, taskStatusOpen, time.Now(), time.Now(), task.AuthorId)

	if err != nil {
		log.Printf("[saveTask] error while saving tasks: %v", err)
		return nil, err
	}

	// TODO add handle events by pipe not right here
	_, err = s.db.Exec(`INSERT into tasks_events (task_id, event_type, occurred_on) values ($1, $2, $3)`, id, "task_created", time.Now())

	if err != nil {
		log.Println(err)
	}

	event := WsEvent{
		Type:  "task_added",
		Event: task,
	}
	wsEventJson, _ := json.Marshal(event)
	hub.Broadcast <- wsEventJson

	return r, nil
}
