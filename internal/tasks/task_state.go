package tasks

type State struct {
	TaskId         string `json:"task_id"`
	DiscussionId   string
	TaskTitle      string `json:"task_title"`
	UnreadComments int    `json:"unread_comments"`
}

type States map[string]*State

func (s States) DiscussionsIds() []string {
	var ids []string

	for _, v := range s {
		ids = append(ids, v.DiscussionId)
	}

	return ids
}

func LoadTaskStateList(ts TaskStorage, userId int, teamId int) map[string]*State {
	states := ts.LoadStates(teamId)
	unreadCommentsAmount := ts.LoadUnreadCommentsAmount(userId, states.DiscussionsIds())

	for taskId, unreadComments := range unreadCommentsAmount {
		states[taskId].UnreadComments = unreadComments
	}

	return states
}
