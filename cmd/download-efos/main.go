package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/net/html/charset"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type EFOSDefinitivo struct {
	RFC      string
	Nombre   string
	FechaSat string
	FechaDof string
	DateSat  time.Time
	DateDof  *time.Time
}

type Service struct {
	db *sqlx.DB
}

func main() {
	fmt.Println("Starting EFOS download")
	downloadDefinitivos()

	fmt.Println("reading csv file")
	efos := readCsvLine()

	fmt.Println("inserting efos to database")

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
	service := Service{db: db}

	service.insert(efos)

	fmt.Println("done")
}

func (s Service) insert(efos []EFOSDefinitivo) {
	sqlQuery := `
	INSERT INTO listado_definitivo 
		(
			rfc, 
			nombre, 
			fecha_publicacion_sat_definitivos_text, 
			fecha_publicacion_dof_definitivos_text,
			fecha_publicacion_sat_definitivos,
			fecha_publicacion_dof_definitivos
		) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	ON CONFLICT (rfc) 
	DO UPDATE SET
		fecha_publicacion_sat_definitivos_text = EXCLUDED.fecha_publicacion_sat_definitivos_text,
		fecha_publicacion_dof_definitivos_text = EXCLUDED.fecha_publicacion_dof_definitivos_text,
		fecha_publicacion_sat_definitivos = EXCLUDED.fecha_publicacion_sat_definitivos,
		fecha_publicacion_dof_definitivos = EXCLUDED.fecha_publicacion_dof_definitivos
	`
	for _, efo := range efos {
		_, err := s.db.Exec(sqlQuery,
			efo.RFC,
			efo.Nombre,
			efo.FechaSat,
			efo.FechaDof,
			efo.DateSat,
			efo.DateDof,
		)
		if err != nil {
			fmt.Printf("there was an error inserting %v: %+v\n", efo.RFC, err)
		}
		fmt.Println("Inserted", efo.RFC)
	}
}

func downloadDefinitivos() {
	if _, err := os.Stat("./definitivos.csv"); err == nil {
		fmt.Println("removing previous file")
		err := os.Remove("./definitivos.csv")
		if err != nil {
			log.Fatalf("error removing definitivos file: %v", err)
		}
	}

	resp, err := http.Get("http://omawww.sat.gob.mx/cifras_sat/Documents/Definitivos.csv")

	if err != nil {
		log.Fatalf("error downloading definitivos file: %v", err)
	}

	defer resp.Body.Close()

	fmt.Println("Creating new definitivos.csv file")

	out, err := os.Create("./definitivos.csv")

	if err != nil {
		log.Fatalf("error creating file definitivos.csv: %v", err)
	}

	defer out.Close()

	red, err := charset.NewReader(resp.Body, "latin1")

	if err != nil {
		log.Fatalf("error creating reader: %v", err)
	}

	_, err = io.Copy(out, red)

	if err != nil {
		log.Fatalf("error copying buffer to file: %v", err)
	}
}

func readCsvLine() []EFOSDefinitivo {
	f, err := os.Open("./definitivos.csv")

	if err != nil {
		log.Fatalf("could not read file definitivos: %v", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	var definitivos []EFOSDefinitivo

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("an error happened while reading csv file: %v", err)
		}

		firstValue := rec[0]
		if strings.HasPrefix(firstValue, "Información actualizada") {
			continue
		}
		if strings.HasPrefix(firstValue, "Listado completo de") {
			continue
		}

		if strings.HasPrefix(firstValue, "No.") {
			continue
		}

		rfc := rec[1]

		if strings.HasPrefix(rfc, "XXXXXXX") {
			continue
		}

		nombre := strings.TrimSpace(rec[2])
		publicacionSatDefinitivos := rec[13]
		publicacionDOFDefinitivos := rec[15]

		if publicacionSatDefinitivos != "" {
			publicacionSatDefinitivos = publicacionSatDefinitivos[:10]
		}

		//"2006-01-02 03:04:05"
		if publicacionDOFDefinitivos != "" {
			publicacionDOFDefinitivos = publicacionDOFDefinitivos[:10]
		}

		var satDate time.Time
		var dofDate *time.Time

		satDate, err = time.ParseInLocation("02/01/2006", publicacionSatDefinitivos, time.FixedZone("UTC", -6*60*60))
		if err != nil {
			fmt.Println("could not parse SAT date", publicacionSatDefinitivos)
		}

		// It could be published by SAT but not yet on DOF
		parsedDofDate, err := time.ParseInLocation("02/01/2006", publicacionDOFDefinitivos, time.FixedZone("UTC", -6*60*60))
		if err != nil {
			fmt.Println("could not parse SAT date", publicacionDOFDefinitivos)
		}

		if parsedDofDate.IsZero() {
			dofDate = nil
		} else {
			dofDate = &parsedDofDate
		}

		efosDefinitivo := EFOSDefinitivo{
			RFC:      rfc,
			Nombre:   nombre,
			FechaSat: publicacionSatDefinitivos,
			FechaDof: publicacionDOFDefinitivos,
			DateSat:  satDate,
			DateDof:  dofDate,
		}

		definitivos = append(definitivos, efosDefinitivo)
	}

	return definitivos
}
