package storage

// TODO move it to platform after resolve issue with WsEvent

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"tasks17-server/internal/tasks"
	"time"
)

type Storage struct {
	*sql.DB
}

func (s *Storage) CreateDiscussion(id string) {
	_, err := s.Exec(`INSERT into discussions (id) values ($1)`, id)

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}
}

func New(db *sql.DB) Storage {
	return Storage{db}
}

func (s *Storage) CompleteSubTask(subTask int) {
	_, _ = s.Exec("UPDATE sub_tasks SET status = $1, closed_at = $2 WHERE id = $3", 1, time.Now(), subTask)
}

func (s *Storage) SaveCommentEvent(comment tasks.Comment) {
	commentJson, err := json.Marshal(comment)
	_, err = s.Exec(`INSERT into tasks_events (task_id, event_type, payload, occurred_on) values ($1, $2, $3, $4)`, comment.DiscussionId, "task_comment_left", commentJson, time.Now())

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}

}

func (s *Storage) SaveComment(comment tasks.Comment) (id int, err error) {
	err = s.QueryRow("INSERT INTO comments (author_id, message, created_at,discussion_id) values ($1,$2,$3,$4) RETURNING id", comment.Author, comment.Message, time.Now(), comment.DiscussionId).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting task comment %v", err)
		return 0, err
	}

	return id, nil
}

func (s *Storage) LoadComments(discussionId string) []tasks.Comment {
	taskCommentsList := make([]tasks.Comment, 0)

	rows, err := s.Query("SELECT comments.id, discussion_id, message, author_id, users.name, created_at FROM comments LEFT JOIN users ON users.id = comments.author_id WHERE discussion_id= $1 ORDER BY created_at", discussionId)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var taskComment tasks.Comment
		err = rows.Scan(&taskComment.Id, &taskComment.DiscussionId, &taskComment.Message, &taskComment.Author, &taskComment.AuthorName, &taskComment.CreatedAt)

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

func (s *Storage) UpdateLastWatchedComment(userId int, discussionId string, commentId int) {
	_, err := s.Exec(`INSERT into discussion_watched_comment (user_id, discussion_id, last_comment_id) values ($1, $2, $3) ON CONFLICT (user_id, discussion_id) DO UPDATE SET last_comment_id = $3`, userId, discussionId, commentId)

	if err != nil {
		log.Printf("Error while inserting task event %v", err)
	}
}

func (s *Storage) LoadSubTasks(taskId string) []tasks.SubTask {
	subTasksList := make([]tasks.SubTask, 0)

	rows, err := s.Query("SELECT id, task_id, stage_id, author_id, status, name, created_at, closed_at FROM sub_tasks WHERE task_id = $1 order by rank", taskId)
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

func (s *Storage) LoadStates(teamId int) tasks.States {
	rows, err := s.Query("SELECT id, title, discussion_id FROM tasks where status = 'open' and team_id = $1", teamId)
	states := make(map[string]*tasks.State)

	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var s tasks.State
		err = rows.Scan(&s.TaskId, &s.TaskTitle, &s.DiscussionId)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			states[s.TaskId] = &s
		}
	}

	return states
}

func (s *Storage) LoadUnreadCommentsAmount(uid int, discussionsIds []string) map[string]int {
	rows, err := s.Query(`SELECT dwc.discussion_id, count(*)
FROM comments
         LEFT JOIN discussion_watched_comment dwc on comments.discussion_id = dwc.discussion_id
where comments.discussion_id = ANY ($1::char(36)[])
  and dwc.user_id = $2
  and comments.id > dwc.last_comment_id
GROUP BY dwc.discussion_id`, pq.Array(discussionsIds), uid)
	defer rows.Close()

	if err != nil {
		log.Panic(err)
	}

	result := map[string]int{}

	for rows.Next() {
		var taskId string
		var unreadComments int

		err = rows.Scan(&taskId, &unreadComments)

		if err != nil {
			log.Println("Error while scanning entity", err)
		} else {
			result[taskId] = unreadComments
		}
	}

	return result
}

func (s *Storage) LoadTask(taskId string) tasks.Task {
	row := s.QueryRow("SELECT id, title, description, status, created_at, author_id, team_id, discussion_id, updated_at FROM tasks WHERE id = $1", taskId)

	var task tasks.Task

	err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.AuthorId, &task.TeamId, &task.DiscussionId, &task.UpdatedAt)

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

	_, err = s.Exec(`INSERT into tasks (id, team_id, discussion_id, title, description, status, created_at, updated_at, author_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, task.Id, task.TeamId, task.DiscussionId, task.Title, task.Description, "open", task.CreatedAt, task.UpdatedAt, task.AuthorId)

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
	err = s.QueryRow(`INSERT into sub_tasks (task_id, stage_id, author_id, rank, status, name, created_at) values ($1, $2, $3, (SELECT count(rank) + 1 FROM sub_tasks WHERE task_id = $1), $4, $5, $6) RETURNING id`, st.TaskId, st.StageId, st.AuthorId, 0, st.Name, time.Now()).Scan(&id)

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
