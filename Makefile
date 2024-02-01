postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root private_tutor_db

dropdb:
	docker exec -it postgres16 dropdb private_tutor_db

.PHONY: postgrs createdb  dropdb