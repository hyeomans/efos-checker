package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrateDownFlag := flag.Bool("down", false, "revert all migrations")
	flag.Parse()

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresSslmode := os.Getenv("POSTGRES_SSLMODE")

	dbConnection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDb, postgresSslmode)

	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		dbConnection,
	)

	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	if *migrateDownFlag {
		fmt.Println("Rolling back migrations")
		err = m.Down()
	} else {
		fmt.Println("Running migrations")
		err = m.Up()
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("\n> Migrations are up to date")
			os.Exit(0)
		}
		log.Fatalf("there was an error running the migrations: %v", err)
	}

	if *migrateDownFlag {
		fmt.Println("Reverted all migrations")
	} else {
		fmt.Println("Ran all migrations")
	}
}
