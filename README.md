Project structure follow rules from https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

## Install

### Install Go

version go 1.13

### Install go-migrate
 https://github.com/golang-migrate/migrate
####Macos
 brew install golang-migrate 
 
#### Ubuntu
https://github.com/golang/go/wiki/Ubuntu

`go mod download`

`docker-compose up -d`

load dump from file to mysql

`go build ./cmd/servid/servid.go`

`export POSTGRESQL_URL=postgres://postgres:tasks17@localhost/tasks17?sslmode=disable`

`migrate -database ${POSTGRESQL_URL} -path database/migrations up `

`go run cmd/servi/servi.go setup`
This command will create organisation, team, and 2 users with credentials:
kir@gangsterelephant.io:kir
serg@gangsterelephant.io:serg

`./servid`
