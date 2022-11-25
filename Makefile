DB_URL=postgresql://root:budgetapidb@localhost:5432/budgetapidb?sslmode=disable

newPostgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=budgetapidb -d postgres:latest

postgres:
	docker start postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root budgetapidb

dropdb:
	docker exec -it postgres dropdb budgetapidb

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/LeandroEstevez/budgetAppAPI/db/sqlc Store

.PHONY: network newPostgres postgres createdb dropdb migrateup migratedown sqlc server mock