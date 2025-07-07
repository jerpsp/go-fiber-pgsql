start:
	docker compose up -d

stop:
	docker compose down

restart:
	make stop && make start

logs:
	docker logs -f go-fiber-api

reset:
	make drop && make create && make migrate && make seed

drop:
	docker exec go-fiber-api go run cmd/cli/main.go dbDrop

create:
	docker exec go-fiber-api go run cmd/cli/main.go dbCreate

migrate:
	docker exec go-fiber-api go run cmd/cli/main.go dbAutoMigrate

seed:
	docker exec go-fiber-api go run cmd/cli/main.go dbSeed	