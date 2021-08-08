migrateup:
	migrate -path db/migration -database "postgresql://root:123456@dbserver:5432/simple_bank?sslmode=disable" -verbose up

migrateupAction:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@dbserver:5432/simple_bank?sslmode=disable" -verbose down

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY:	migrateup migratedown migratecreate sqlc test migrateupAction
