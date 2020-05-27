package tasks

type State struct {
	TaskId         string `json:"task_id"`
	TaskTitle      string `json:"task_title"`
	UnreadComments int    `json:"unread_comments"`
}

func LoadTaskStateList(ts TaskStorage, userId int, teamId int) map[string]*State {
	unreadCommentsAmount := ts.LoadUnreadCommentsAmount(userId, teamId)
	states := ts.LoadStates(teamId)

	for taskId, unreadComments := range unreadCommentsAmount {
		states[taskId].UnreadComments = unreadComments
	}

	return states
}
