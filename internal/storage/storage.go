package storage

// TODO move it to platform after resolve issue with WsEvent

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"tasks17-server/internal/tasks"
	"time"
)

type Storage struct {
	*sql.DB
}

func New(db *sql.DB) Storage {
	return Storage{db}
}

func (s *Storage) CompleteSubTask(subTask int) {
	_, _ = s.Exec("UPDATE sub_tasks SET status = $1, closed_at = $2 WHERE id = $3", 1, time.Now(), subTask)
}

func (s *Storage) SaveCommentEvent(comment tasks.TaskComment) {
	commentJson, err := json.Marshal(comment)
	_, err = s.Exec(`INSERT into tasks_events (task_id, event_type, payload, occurred_on) values ($1, $2, $3, $4)`, comment.TaskId, "task_comment_left", commentJson, time.Now())

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}

}

func (s *Storage) SaveComment(comment tasks.TaskComment) (id string, err error) {
	generatedId, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	comment.CreatedAt = time.Now()
	comment.Id = generatedId.String()

	_, err = s.Exec("INSERT INTO task_comments (id, author_id, message, created_at, task_id) values ($1,$2,$3,$4,$5)", comment.Id, comment.Author, comment.Message, comment.CreatedAt, comment.TaskId)

	if err != nil {
		log.Printf("Error while inserting task comment %v", err)
		return "", err
	}

	return comment.Id, nil
}

func (s *Storage) LoadComments(taskId string) []tasks.TaskComment {
	taskCommentsList := make([]tasks.TaskComment, 0)

	rows, err := s.Query("SELECT task_comments.id, task_id, message, author_id, users.name, created_at FROM task_comments LEFT JOIN users ON users.id = task_comments.author_id WHERE task_id= $1 ORDER BY created_at", taskId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var taskComment tasks.TaskComment
		err = rows.Scan(&taskComment.Id, &taskComment.TaskId, &taskComment.Message, &taskComment.Author, &taskComment.AuthorName, &taskComment.CreatedAt)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			taskCommentsList = append(taskCommentsList, taskComment)
		}
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return taskCommentsList
}

func (s *Storage) FindLastEventByCommentId(id string) int {
	rows, err := s.Query("SELECT id FROM tasks_events WHERE event_type = 'task_comment_left' AND payload LIKE $1 LIMIT 1", "%"+id+"%")
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var eventId int
		err = rows.Scan(&eventId)

		if err != nil {
			log.Println("Error while scanning entity", err)
		}
		return eventId
	}

	return 0

}

func (s *Storage) UpdateLastWatchedComment(userId int, taskId string, commentId string) {
	// TODO return error if it occurrenced
	id := s.FindLastEventByCommentId(commentId)
	_, err := s.Exec(`INSERT into task_last_watched_event (user_id, task_id, last_event_id) values ($1, $2, $3) ON CONFLICT (user_id, task_id) DO UPDATE SET last_event_id = $3`, userId, taskId, id)

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}
}

func (s *Storage) LoadSubTasks(taskId string) []tasks.SubTask {
	subTasksList := make([]tasks.SubTask, 0)

	rows, err := s.Query("SELECT id, task_id, stage_id, author_id, status, name, created_at, closed_at FROM sub_tasks WHERE task_id = $1", taskId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var st tasks.SubTask
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

func (s *Storage) LoadTaskStages(taskId string) []tasks.Stage {
	list := make([]tasks.Stage, 0)

	rows, err := s.Query("SELECT DISTINCT sts.id, sts.name, sts.rank FROM sub_tasks LEFT JOIN sub_task_stages sts on sub_tasks.stage_id = sts.id WHERE task_id = $1", taskId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var s tasks.Stage

		err = rows.Scan(&s.Id, &s.Name, &s.Rank)

		if err != nil {
			log.Println("Error while scanning Stage entity", err)
		} else {
			list = append(list, s)
		}
	}

	if err = rows.Err(); err != nil {
		log.Panic(err)
	}

	return list
}

func (s *Storage) LoadStages(teamId int) []tasks.Stage {
	list := make([]tasks.Stage, 0)

	rows, err := s.Query("SELECT sts.id, sts.name, sts.rank FROM sub_task_stages AS sts WHERE team_id = $1", teamId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var s tasks.Stage

		err = rows.Scan(&s.Id, &s.Name, &s.Rank)

		if err != nil {
			log.Println("Error while scanning Stage entity", err)
		} else {
			list = append(list, s)
		}
	}

	if err = rows.Err(); err != nil {
		log.Panic(err)
	}

	return list
}

func (s *Storage) CloseTask(id string) error {
	r, err := s.Exec("UPDATE tasks SET status = $1 WHERE id = $2", "closed", id)

	if err != nil {
		return err
	}

	if affected, _ := r.RowsAffected(); affected != 1 {
		return errors.New(fmt.Sprintf("task with id: %s not found", id))
	}

	return nil
}

func (s *Storage) UpdateDescription(taskId string, description string) error {
	r, err := s.Exec("UPDATE tasks SET description = $1 WHERE id = $2", description, taskId)

	if err != nil {
		return err
	}

	if affected, _ := r.RowsAffected(); affected != 1 {
		return errors.New(fmt.Sprintf("task with id: %s not found", taskId))
	}

	return nil
}

func (s *Storage) LoadStates(teamId int) map[string]*tasks.State {
	rows, err := s.Query("SELECT id, title FROM tasks where status = 'open' and team_id = $1", teamId)
	states := make(map[string]*tasks.State)

	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var s tasks.State
		err = rows.Scan(&s.TaskId, &s.TaskTitle)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			states[s.TaskId] = &s
		}
	}

	return states
}

func (s *Storage) LoadEvents(userId int, teamId int) []tasks.Event {
	events := make([]tasks.Event, 0)
	rows, err := s.Query(`SELECT te.task_id, te.event_type, te.payload, te.occurred_on 
										FROM tasks_events te
											LEFT JOIN tasks ON te.task_id = tasks.id
											LEFT JOIN task_last_watched_event tlwe ON te.task_id = tlwe.task_id AND tlwe.user_id = $1
									WHERE team_id = $2 AND (te.id > tlwe.last_event_id or tlwe.last_event_id IS NULL) and status = 'open'`, userId, teamId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var payload sql.NullString
		var e tasks.Event

		err = rows.Scan(&e.TaskId, &e.EventType, &payload, &e.OccurredOn)
		e.Payload = payload.String

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			events = append(events, e)
		}
	}

	return events
}

func (s *Storage) LoadTask(taskId string) tasks.Task {
	row := s.QueryRow("SELECT id, title, description, status, created_at, author_id, team_id, updated_at FROM tasks WHERE id = $1", taskId)

	var task tasks.Task

	err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.AuthorId, &task.TeamId, &task.UpdatedAt)

	if err != nil {
		log.Println("Error while scanning entity", err)
	}

	return task
}

func (s *Storage) LoadTasks(teamId int) []tasks.Task {
	taskList := make([]tasks.Task, 0)

	rows, err := s.Query("SELECT id, title, description, status, created_at FROM tasks WHERE team_id=$1", teamId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var task tasks.Task
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

func (s *Storage) SaveTask(task *tasks.Task) (err error) {
	task.PreSave()

	_, err = s.Exec(`INSERT into tasks (id, team_id, title, description, status, created_at, updated_at, author_id) values ($1, $2, $3, $4, $5, $6, $7, $8)`, task.Id, task.TeamId, task.Title, task.Description, "open", task.CreatedAt, task.UpdatedAt, task.AuthorId)

	if err != nil {
		log.Printf("[SaveTask] error while saving tasks: %v", err)
		return err
	}

	// TODO add handle events by pipe not right here
	_, err = s.Exec(`INSERT into tasks_events (task_id, event_type, occurred_on) values ($1, $2, $3)`, task.Id, "task_created", time.Now())

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (s *Storage) SaveSubTask(st tasks.SubTask) (id int, err error) {
	err = s.QueryRow(`INSERT into sub_tasks (task_id, stage_id, author_id, status, name, created_at) values ($1, $2, $3, $4, $5, $6) RETURNING id`, st.TaskId, st.StageId, st.AuthorId, 0, st.Name, time.Now()).Scan(&id)

	if err != nil {
		log.Printf("[addSubTask] error while adding sub-tasks: %v", err)
		return 0, err
	}

	// TODO add handle events by pipe not right here
	_, err = s.Exec(`INSERT into tasks_events (task_id, event_type, occurred_on) values ($1, $2, $3)`, st.TaskId, "sub_task_added", time.Now())

	if err != nil {
		log.Println("Error while task event inserting", err)
	}

	return id, nil
}

func (s *Storage) LoadTeams(orgId int) tasks.Teams {
	list := make([]tasks.Team, 0)

	rows, err := s.Query(`SELECT id, organisation_id,  name FROM teams WHERE organisation_id = $1`, orgId)

	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var team tasks.Team
		err = rows.Scan(&team.Id, &team.OrgId, &team.Name)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			list = append(list, team)
		}
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return list
}

func (s *Storage) RemoveTask(id string) {
	// TODO add error processing here
	_, err := s.Exec(`DELETE FROM tasks where id = $1`, id)

	if err != nil {
		log.Println(err)
	}
}

func (s *Storage) SaveTeam(team tasks.Team) (id int, err error) {
	err = s.QueryRow(`INSERT into teams (organisation_id, name) values ($1, $2) RETURNING id`, team.OrgId, team.Name).Scan(&id)

	if err != nil {
		log.Printf("[addSubTask] error while adding sub-tasks: %v", err)
		return 0, err
	}

	return id, err
}

func (s *Storage) RemoveTeam(id int) {
	// TODO add error processing here
	_, err := s.Exec(`DELETE FROM teams where id = $1`, id)

	if err != nil {
		log.Println(err)
	}
}

func (s *Storage) LoadTeam(id int) tasks.Team {
	row := s.QueryRow("SELECT id, organisation_id, name FROM teams WHERE id = $1", id)

	var team tasks.Team

	err := row.Scan(&team.Id, &team.OrgId, &team.Name)

	if err != nil {
		log.Println("Error while scanning entity", err)
	}

	return team
}

func (s *Storage) SaveStage(stage tasks.Stage) (id int, err error) {
	err = s.QueryRow(`INSERT into sub_task_stages (team_id, rank, name) values ($1, $2, $3) RETURNING id`, stage.TeamId, stage.Rank, stage.Name).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting stage %v", err)
		return 0, err
	}

	return id, err
}

func (s *Storage) SaveOrg(org tasks.Organisation) (id int, err error) {
	err = s.QueryRow(`INSERT into organisations (name) values ($1) RETURNING id`, org.Name).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting org %v", err)
		return 0, err
	}

	return id, err
}

func (s *Storage) RemoveOrg(id int) {
	_, err := s.Exec(`DELETE FROM organisations where id = $1`, id)

	if err != nil {
		log.Println(err)
	}
}

func (s *Storage) SaveUser(member tasks.User) (id int, err error) {
	err = s.QueryRow(`INSERT into users (organisation_id, name, email, password) values ($1, $2, $3, $4) RETURNING id`, member.OrgId, member.Name, member.Email, member.Password).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting org %v", err)
		return 0, err
	}

	return id, err
}
