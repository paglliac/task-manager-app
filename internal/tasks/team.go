package tasks

type TeamModule struct {
	Storage TaskStorage
}

func (t TeamModule) LoadTeams(orgId int) []Team {
	return t.Storage.LoadTeams(orgId)
}

type Team struct {
	Id    int    `json:"id"`
	OrgId int    `json:"org_id"`
	Name  string `json:"name"`
}

type Teams []Team

func (teams Teams) ById(id int) *Team {
	for _, team := range teams {
		if team.Id == id {
			return &team
		}
	}

	return nil
}
