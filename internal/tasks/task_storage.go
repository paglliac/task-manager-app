package tasks

type TaskStorage interface {
	LoadTasks(teamId int) []Task
	// TODO move teams outside task storage
	LoadTask(id string) Task
	SaveTask(t *Task) error
	CloseTask(id string) error
	LoadStates(teamId int) map[string]*State
	LoadUnreadCommentsAmount(uid int, teamId int) map[string]int
	UpdateDescription(id string, description string) error

	LoadTeams(orgId int) Teams
	LoadTeam(id int) Team

	LoadStages(teamId int) []Stage
	LoadTaskStages(taskId string) []Stage
	LoadSubTasks(taskId string) []SubTask
	CompleteSubTask(subTaskId int)
	SaveSubTask(st SubTask) (int, error)

	SaveComment(comment TaskComment) (id int, err error)
	SaveCommentEvent(comment TaskComment)
	UpdateLastWatchedComment(userId int, taskId string, commentId int)
	LoadComments(taskId string) []TaskComment
}
