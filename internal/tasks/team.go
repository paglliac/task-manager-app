package tasks

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
