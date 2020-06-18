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
	Id          int    `json:"id"`
	ProjectId   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}
