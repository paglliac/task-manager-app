package tasks_test

import (
	"os"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/storage"
	"tasks17-server/internal/tasks"
	"testing"
)

type hub struct {
	events []platform.WsEvent
}

func (h *hub) Handle(event platform.WsEvent) {
	h.events = append(h.events, event)
}

var s storage.Storage
var h hub

func TestMain(m *testing.M) {
	dbHost := os.Getenv("DB_HOST")
	db, _ := platform.InitDb(dbHost)

	s = storage.New(db)

	h = hub{}
	tasks.Init(&h)

	code := m.Run()
	os.Exit(code)
}
