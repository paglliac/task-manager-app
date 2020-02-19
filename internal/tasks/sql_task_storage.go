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

func (s *SqlTaskStorage) getDb() *platform.Storage {
	return &s.db
}

func (s *SqlTaskStorage) loadTasksState(userId int) (map[string]*TaskState, []TaskEvent) {
	states := make(map[string]*TaskState)
	events := make([]TaskEvent, 0)
	rows, err := s.db.Query(`SELECT t.id, te.id, t.title, te.event_type, te.payload, te.occurred_on
						FROM tasks t
								 LEFT JOIN task_last_watched_event tlwe ON t.id = tlwe.task_id
								 RIGHT JOIN tasks_events te ON t.id = te.task_id
						WHERE 
						      (user_id = ? OR user_id IS NULL)
						  AND (te.id > tlwe.last_event_id OR tlwe.last_event_id IS NULL OR te.event_type != 'task_comment_left')`, userId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {

		var (
			eventType  string
			taskId     string
			eventId    int
			title      string
			payload    sql.NullString
			occurredOn time.Time
		)

		err = rows.Scan(&taskId, &eventId, &title, &eventType, &payload, &occurredOn)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			if _, ok := states[taskId]; !ok {
				var s TaskState
				s.TaskId = taskId
				s.TaskName = title

				states[taskId] = &s
			}

			var e TaskEvent

			e.taskId = taskId
			e.occurredOn = occurredOn
			e.eventType = eventType
			e.payload = payload.String

			events = append(events, e)
		}
	}

	return states, events
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
