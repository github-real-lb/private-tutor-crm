postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

start:
	docker start postgres16

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres tutor_db

dropdb:
	docker exec -it postgres16 dropdb tutor_db

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/tutor_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/tutor_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./... -cover

server:
	go run .

mockery:
	mockery --output db/mocks --filename store.go --outpkg mocks  --dir db/sqlc --name Store --structname MockStore

.PHONY: postgrs start createdb  dropdb migrateup migratedown sqlc test server mockery

localcreatedb:
	createdb --username=postgres --password=secret tutor_db

localdropdb:
	dropdb --username=postgres --password=secret tutor_db

localmigrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/tutor_db?sslmode=disable" -verbose up

localmigratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/tutor_db?sslmode=disable" -verbose down
