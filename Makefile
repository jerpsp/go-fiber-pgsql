start:
	docker compose up -d

stop:
	docker compose down

restart:
	make stop && make start

logs:
	docker logs -f go-fiber-api

rebuild:
	docker build -t go-fiber-api .

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

test:
	go test ./internal/api/... -v -covermode count -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html

prod-tag:
	git tag -d release_$(release) || true
	git tag release_$(release)
	git push origin release_$(release)