package tasks

import (
	"tasks17-server/internal/platform"
	"time"
)

type EventHandler interface {
	Handle(event platform.WsEvent)
}

type Event struct {
	TaskId     string
	EventType  string
	OccurredOn time.Time
	Payload    string
}
