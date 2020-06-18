package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/storage"
	"tasks17-server/internal/tasks"
	"testing"
)

type dbSetup struct {
	orgId    int
	teamId   int
	teamName string
	user     user
	project  tasks.Project
}

type user struct {
	id       int
	email    string
	password string
	teamId   int
}

var r http.Handler
var s storage.Storage

var setup = dbSetup{}

var authenticator auth.Authenticator

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	dbHost := os.Getenv("DB_HOST")

	db, err := platform.InitDb(dbHost)
	if err != nil {
		panic(err)
	}

	s = storage.New(db)

	hub := platform.InitHub()

	tasks.Init(hub)
	authenticator = auth.NewAuthenticator(db)

	r = CreateRouter(hub, db)
	setUpDatabase()
	defer tearDownDatabase()
	m.Run()
}

func tearDownDatabase() {
	removeOrg()
}

func setUpDatabase() {
	createOrg()
	createTeam()
	createTeamMember()
	createProject()
}

func authRequest(req *http.Request, user *user) {
	token, _ := authenticator.IssueToken(user.email, user.password)
	req.Header.Set("Authorization", token)
}

func toJsonReader(o interface{}) *bytes.Reader {
	res, _ := json.Marshal(o)
	return bytes.NewReader(res)
}

type jsonResponse struct {
	parsedResponse map[string]interface{}
	rawResponse    *http.Response
	parsed         bool
}

func (r *jsonResponse) parse() {
	if r.parsed {
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.rawResponse.Body)
	r.rawResponse.Body.Close() //  must close
	r.rawResponse.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.rawResponse.Body).Decode(&r.parsedResponse); err != nil {
		bodyString := string(bodyBytes)
		log.Fatalf("Error while decode json response: %v. %s", err, bodyString)
	}

	r.parsed = true
}

func (r *jsonResponse) getRaw(key string) interface{} {
	r.parse()

	v := r.parsedResponse[key]

	return v
}

func (r *jsonResponse) get(key string) string {
	r.parse()

	v := r.parsedResponse[key].(string)
	return v
}

type Requester struct {
	user       *user
	test       *testing.T
	skipErrors bool
}

func NewRequester(t *testing.T) *Requester {
	rr := Requester{
		test:       t,
		skipErrors: false,
	}
	return &rr
}

func (rr *Requester) auth(u *user) *Requester {
	rr.user = u

	return rr
}

func (rr *Requester) get(target string) jsonResponse {
	req := httptest.NewRequest("GET", target, nil)

	return rr.processRequest(req)
}

func (rr *Requester) post(target string, body interface{}) jsonResponse {
	var req *http.Request

	if body == nil {
		req = httptest.NewRequest("POST", target, nil)
	} else {
		req = httptest.NewRequest("POST", target, toJsonReader(body))
	}

	return rr.processRequest(req)
}

func (rr *Requester) processRequest(req *http.Request) jsonResponse {
	if rr.user != nil {
		authRequest(req, rr.user)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if rr.skipErrors {
		return jsonResponse{rawResponse: w.Result()}
	}

	if w.Result().StatusCode != http.StatusOK {
		rr.test.Errorf("Response should be success, got %d", w.Result().StatusCode)
		rr.test.FailNow()
	}

	return jsonResponse{rawResponse: w.Result()}

}

func createOrg() {
	id, _ := s.SaveOrg(tasks.Organisation{
		Name: "Test org",
	})

	setup.orgId = id
}

func removeOrg() {
	s.RemoveOrg(setup.orgId)
}

func createTeam() {
	id, _ := s.SaveTeam(tasks.Team{
		OrgId: setup.orgId,
		Name:  "Test team",
	})

	setup.teamId = id
	setup.teamName = "Test team"
}

func createTeamMember() {
	u := tasks.User{
		OrgId:    setup.orgId,
		Name:     "User",
		Email:    "user@example.com",
		Password: "password",
	}

	id, err := s.SaveUser(u)

	if err != nil {
		panic(err)
	}

	setup.user = user{
		id:       id,
		email:    u.Email,
		password: u.Password,
		teamId:   setup.teamId,
	}

}

func createProject() {
	discussionId := uuid.New().String()
	s.CreateDiscussion(discussionId)
	p := tasks.Project{
		OrgId:        setup.orgId,
		Name:         "Project A",
		Description:  "Project A Discussion",
		Status:       0,
		DiscussionId: discussionId,
	}
	id := s.CreateProject(p)
	p.Id = id
	setup.project = p

}

func (u user) createTask() (tasks.Task, context.CancelFunc) {
	task := tasks.Task{
		TeamId:      u.teamId,
		Title:       "New Task",
		Description: "New Task Description",
		AuthorId:    u.id,
	}

	_, err := tasks.CreateTask(&s, &task)

	if err != nil {
		panic(err)
	}

	return task, func() {
		s.RemoveTask(task.Id)
	}
}
