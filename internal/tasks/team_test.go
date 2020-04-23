package tasks_test

import (
	"tasks17-server/internal/tasks"
	"testing"
)

func TestCreateTeam(t *testing.T) {
	// Task creation
	expectedTeam := tasks.Team{
		OrgId: 1,
		Name:  "Team A",
	}

	id, err := s.SaveTeam(expectedTeam)
	defer s.RemoveTeam(id)

	if err != nil {
		t.Error(err)
	}
	expectedTeam.Id = id

	// Check task has been created
	loadedTeam := s.LoadTeam(id)

	if loadedTeam != expectedTeam {
		t.Errorf("Task not created correctly \n %+v \n %+v", expectedTeam, loadedTeam)
	}

	loadedTeams := s.LoadTeams(1)

	t2 := loadedTeams.ById(expectedTeam.Id)

	if t2 == nil {
		t.Errorf("Team not found in team list")
		return
	}

	if expectedTeam != *t2 {
		t.Errorf("Team not created correctly \n %+v \n %+v", expectedTeam, t2)
	}
}
