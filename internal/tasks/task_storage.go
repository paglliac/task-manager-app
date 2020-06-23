package tasks

type TaskStorage interface {
	LoadTasks(teamId int) []Task
	LoadTasksByProjectStages(ids []int) []Task
	// TODO move teams outside task storage
	LoadTask(id string) Task
	SaveTask(t *Task) error
	CloseTask(id string) error
	LoadStates(teamId int) States
	LoadUnreadCommentsAmount(uid int, discussionsIds []string) map[string]int
	UpdateDescription(id string, description string) error

	LoadTeams(orgId int) Teams
	SaveTeam(team Team) (id int, err error)
	LoadTeam(id int) Team

	LoadStages(teamId int) []Stage
	LoadTaskStages(taskId string) []Stage
	LoadSubTasks(taskId string) []SubTask
	CompleteSubTask(subTaskId int)
	SaveSubTask(st SubTask) (int, error)

	SaveComment(comment Comment) (id int, err error)
	UpdateLastWatchedComment(userId int, taskId string, commentId int)
	LoadComments(discussionId string) []Comment
	CreateDiscussion(id string)

	CreateProject(project Project) int
	LoadProjects(oid int) []Project
	LoadProject(id int) Project
	CreateProjectStage(stage ProjectStage) int
	LoadProjectStages(pid int) []ProjectStage
	UpdateProject(project Project)
}
