DB_URL=postgresql://root:secret@localhost:5432/inventory_system?sslmode=disable

postgress:
# TO GENERATE A NEW POSTGRES CONTAINER
# 1. docker rm -f postgres13
# 2. make postgress
	docker run --name postgres13 --network inventory_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18beta3-alpine

createdb:
# RUN POSTGRES SHELL THROUGH DOCKER
	docker exec -it postgres13 createdb --username=root --owner=root inventory_system

dropdb:
	docker exec -it postgres13 dropdb inventory_system

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1


migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1


startpostgress: 
	docker start postgres13

sqlc:
	sqlc generate

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql

test:
	go test -v --cover ./...

server:
	# go run main.go
	air

.PHONY: server createdb dropdb postgress migrateup migratedown startpostgress sqlc migratedown1 migrateup1 db_docs db_schema
