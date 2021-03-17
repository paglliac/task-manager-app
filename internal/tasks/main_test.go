package tasks_test

import (
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
var hh *platform.Hub

var setup = dbSetup{}

type dbSetup struct {
	orgId    int
	teamId   int
	teamName string
	users    users
}

type users []user

func (u users) main() *user {
	for _, user := range u {
		return &user
	}

	return nil
}

func (u users) secondary() *user {
	for i, user := range u {
		if i == 1 {
			return &user
		}
	}

	return nil
}

type user struct {
	id       int
	email    string
	password string
	teamId   int
}

func TestMain(m *testing.M) {
	config := platform.InitConfig()

	db, _ := platform.InitDb(config.DbConfig)

	s = storage.New(db)

	h = hub{}
	tasks.Init(&h)
	hh = platform.InitHub()

	tearDown := setUpDatabase()
	defer tearDown()

	m.Run()
}

func setUpDatabase() (tearDown func()) {
	createOrg()
	createTeam()
	createTeamMembers()

	return func() {
		removeOrg()
	}
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

func createTeamMembers() {
	u1 := tasks.User{
		OrgId:    setup.orgId,
		Name:     "First User",
		Email:    "first_user@example.com",
		Password: "password",
	}

	id, err := s.SaveUser(u1)

	if err != nil {
		panic(err)
	}

	setup.users = append(setup.users, user{
		id:       id,
		email:    u1.Email,
		password: u1.Password,
		teamId:   setup.teamId,
	})

	u2 := tasks.User{
		OrgId:    setup.orgId,
		Name:     "Second User",
		Email:    "second_user@example.com",
		Password: "password",
	}

	id, err = s.SaveUser(u2)

	if err != nil {
		panic(err)
	}

	setup.users = append(setup.users, user{
		id:       id,
		email:    u2.Email,
		password: u2.Password,
		teamId:   setup.teamId,
	})
}
