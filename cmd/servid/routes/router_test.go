package routes

import (
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/http/httptest"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/tasks"
	"testing"
)

func TestCreateTaskEndpoint(t *testing.T) {
	db, err := platform.InitDB()

	if err != nil {
		panic(err)
	}

	hub := platform.InitHub()

	r := CreateRouter()

	tasks.InitTasksModule(db, hub)

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Error("Response should be success")
	}
}