package tasks

type Project struct {
	Id           int    `json:"id"`
	OrgId        int    `json:"org_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Status       int    `json:"status"`
	DiscussionId string `json:"discussion_id"`
}

type ProjectStage struct {
	id          int
	projectId   int
	name        string
	description string
	status      int
}
