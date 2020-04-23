package tasks

type User struct {
	Id    int
	OrgId int
	Name  string
	Email string

	// TODO store encoded password
	Password string
}
