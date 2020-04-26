package tasks

import "log"

type State struct {
	TaskId         string `json:"task_id"`
	TaskTitle      string `json:"task_title"`
	UnreadComments int    `json:"unread_comments"`
}

func LoadTaskStateList(ts TaskStorage, userId int, teamId int) map[string]*State {
	events := ts.LoadEvents(userId, teamId)
	states := ts.LoadStates(teamId)

	for _, event := range events {
		if event.EventType == "task_comment_left" {
			if _, ok := states[event.TaskId]; ok {
				states[event.TaskId].UnreadComments++
			} else {
				log.Printf("ERR event for not exists task. Event task id: %s", event.TaskId)
			}
		}
	}

	return states
}
