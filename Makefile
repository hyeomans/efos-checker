migrateup:
	go run cmd/migrate/migrate.go

migratedown:
	go run cmd/migrate/migrate.go -down

downloadefos:
	go run cmd/download-efos/main.go

composeup:
	docker-compose up -d

composestop:
	docker-compose stop

searchefos:
	go run cmd/search-efos/main.go

.PHONY: migrateup migratedown downloadefos composeup composestop searchefos