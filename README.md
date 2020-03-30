Project structure follow rules from https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

## Install

### Install Go

version go 1.13

#### Ubuntu
https://github.com/golang/go/wiki/Ubuntu

`go mod download`

`docker-compose up -d`

load dump from file to mysql

`go build ./cmd/servid/servid.go`

`export POSTGRESQL_URL=postgres://postgres:tasks17@localhost/tasks17?sslmode=disable`

`migrate -database ${POSTGRESQL_URL} -path database/migrations up `

`go run cmd/servi/servi.go createUser -username=kir -email=kir@gangsterelephant.io -password=kir`

`./servid`
