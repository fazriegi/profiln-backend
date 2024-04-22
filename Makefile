migrateup:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/profiln" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/profiln" -verbose down
sqlc:
	sqlc generate

.PHONY: migrateup migratedown