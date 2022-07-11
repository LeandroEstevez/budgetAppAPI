postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=budgetdb -d postgres:latest

createdb:
	docker exec -it postgresdb createdb --username=root --owner=root budgetdb

dropdb:
	docker exec -it postgresdb dropdb budgetdb

migrateup:
	migrate -path db/migration -database "postgresql://root:budgetdb@localhost:5432/budgetdb?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:budgetdb@localhost:5432/budgetdb?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc