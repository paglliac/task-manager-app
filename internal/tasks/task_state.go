package tasks

type State struct {
	TaskId         string `json:"task_id"`
	DiscussionId   string `json:"discussion_id"`
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

	for discussionId, unreadComments := range unreadCommentsAmount {
		for i, v := range states {
			if v.DiscussionId == discussionId {
				states[i].UnreadComments = unreadComments
			}
		}
	}

	return states
}
