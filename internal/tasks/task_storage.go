package tasks

type TaskStorage interface {
	LoadTasks(teamId int) []Task
	// TODO move teams outside task storage
	LoadTask(id string) Task
	SaveTask(t *Task) error
	CloseTask(id string) error
	LoadStates(teamId int) map[string]*State
	LoadEvents(uid int, teamId int) []Event
	UpdateDescription(id string, description string) error

	LoadTeams(orgId int) Teams
	LoadTeam(id int) Team

	LoadStages(teamId int) []Stage
	LoadTaskStages(taskId string) []Stage
	LoadSubTasks(taskId string) []SubTask
	CompleteSubTask(subTaskId int)
	SaveSubTask(st SubTask) (int, error)

	SaveComment(comment TaskComment) (id string, err error)
	SaveCommentEvent(comment TaskComment)
	UpdateLastWatchedComment(userId int, taskId string, commentId string)
	FindLastEventByCommentId(id string) int
	LoadComments(taskId string) []TaskComment
}
