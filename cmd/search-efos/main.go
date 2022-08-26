package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Result struct {
	Nombre string
	RFC    string
	Dist   int
}

func main() {
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

	db, err := sqlx.Connect("postgres", dbConnection)

	if err != nil {
		log.Fatalf("could not create db pool: %v", err)
	}

	fmt.Println("reading csv")

	nombres := readNombres()

	fmt.Println("searching in database")

	sqlQuery := `
		SELECT nombre, rfc, $1 <<-> nombre AS dist
		FROM listado_definitivo
		WHERE nombre ilike $2
		LIMIT 1
	`

	for _, nombre := range nombres {
		var r Result
		err := db.Get(&r, sqlQuery, nombre, "%"+nombre+"%")
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}

			fmt.Println(err)
			continue
		}

		fmt.Printf("%+v\n", r)
	}
}

func readNombres() []string {
	f, err := os.Open("./nombres.csv")

	if err != nil {
		log.Fatalf("could not read file definitivos: %v", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	var nombres []string

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("an error happened while reading csv file: %v", err)
		}

		if rec[0] == "" {
			continue
		}
		nombre := rec[0]
		nombres = append(nombres, nombre)
	}

	return nombres
}
