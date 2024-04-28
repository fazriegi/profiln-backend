migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@192.168.100.14:5432/profiln" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@192.168.100.14:5432/profiln" -verbose down
sqlc:
	sqlc generate

.PHONY: migrateup migratedown