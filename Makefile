db:
	docker-compose exec database psql --username=postgres tasks17 -c "drop schema public cascade; create schema public"
	migrate -database postgres://postgres:tasks17@localhost/tasks17?sslmode=disable -path database/migrations up
	go run cmd/servi/servi.go setup