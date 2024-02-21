postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

start:
	docker start postgres16

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root tutor_db

dropdb:
	docker exec -it postgres16 dropdb tutor_db

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/tutor_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/tutor_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./... -cover

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb  github.com/github-real-lb/tutor-management-web/db/sqlc Store

mockery:
	mockery --output db/mocks --filename store.go --outpkg mocks  --dir db/sqlc --name Store --structname MockStore

.PHONY: postgrs start createdb  dropdb migrateup migratedown sqlc test server mock mockery