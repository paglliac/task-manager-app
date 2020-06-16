package tasks

type Project struct {
	id          int
	orgId       int
	name        string
	description string
}

type ProjectStage struct {
	id          int
	projectId   int
	name        string
	description string
}
