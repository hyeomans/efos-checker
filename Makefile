migrateup:
	go run cmd/migrate/migrate.go

migratedown:
	go run cmd/migrate/migrate.go -down

.PHONY: migrateup migratedown