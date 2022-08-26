migrateup:
	go run cmd/migrate/migrate.go

migratedown:
	go run cmd/migrate/migrate.go -down

downloadefos:
	go run cmd/download-efos/main.go

sqlc:
	sqlc generate

.PHONY: migrateup migratedown downloadefos sqlc