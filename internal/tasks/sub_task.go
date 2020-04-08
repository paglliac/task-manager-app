package tasks

import (
	"database/sql"
	"log"
	"time"
)

type SubTask struct {
	Id        int       `json:"id"`
	TaskId    string    `json:"task_id"`
	StageId   int       `json:"stage_id"`
	AuthorId  int       `json:"author_id"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	ClosedAt  time.Time `json:"closed_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Stage struct {
	Id   int
	Name string
	Rank int
}

func LoadSubTasks(taskId string) []SubTask {
	subTasksList := make([]SubTask, 0)

	rows, err := taskStorage.getDb().Query("SELECT id, task, stage, author, status, name, created_at, closed_at FROM sub_tasks WHERE task = $1", taskId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var st SubTask
		var closedAt sql.NullTime

		err = rows.Scan(&st.Id, &st.TaskId, &st.StageId, &st.AuthorId, &st.Status, &st.Name, &st.CreatedAt, &closedAt)
		st.ClosedAt = closedAt.Time

		if err != nil {
			log.Println("Error while scanning SubTask entity", err)
		} else {
			subTasksList = append(subTasksList, st)
		}
	}

	if err = rows.Err(); err != nil {
		log.Panic(err)
	}

	return subTasksList
}
