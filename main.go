package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Place struct {
	Country       string
	City          sql.NullString
	TelephoneCode int `db:"telecode"`
}

func main() {
	var db *sqlx.DB
	var err error

	//Loading environment variables for DATABASE connection
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	//open and connect to the database at the same time
	db, err = sqlx.Connect(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	}

	schema := `CREATE TABLE IF NOT EXISTS place (
				country text,
				city text NULL,
				telecode integer);`

	//execute a query on the server
	result, err2 := db.Exec(schema)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("**************")
	fmt.Println(result)
	fmt.Println("**************")

	// MustExec - panics on error
	cityState := `INSERT INTO place (country, telecode) VALUES ($1, $2);`
	countryCity := `INSERT INTO place (country, city, telecode) VALUES ($1, $2, $3);`
	db.MustExec(cityState, "Hong Kong", 852)
	db.MustExec(cityState, "Singapore", 65)
	db.MustExec(countryCity, "South Africa", "Johannesburg", 27)

	// fetch all places from db
	rows, err := db.Query("SELECT country, city, telecode FROM place")
	if err != nil {
		log.Fatal(err)
	}
	// iterate over each row
	for rows.Next() {
		var place Place
		err = rows.Scan(&place.Country, &place.City, &place.TelephoneCode)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(place)
	}
	// check the error from the rows
	err = rows.Err()

	// fetch a SINGLE ROW
	row := db.QueryRow("SELECT * FROM PLACE where telecode = $1", 852)
	var telecode int
	err = row.Scan(&telecode)
	fmt.Println("Telecode :", telecode)

	p := Place{}
	pp := []Place{}

	// put the first row directly to p
	err = db.Get(&p, "SELECT * FROM place LIMIT 1")
	fmt.Println("Limit 1: ", p)

	// put the places(rows) with telecode>50 into the slice pp
	err = db.Select(&pp, "SELECT * FROM place WHERE telecode > $1", 50)

	fmt.Println("Places with telecom > 50")
	for _, place := range pp {
		fmt.Println(place)
	}

	// they work with regular types as well
	var id int
	err = db.Get(&id, "SELECT count(*) FROM place")
	fmt.Println("Count : ", id)

	// fetch at most 10 places(rows)
	var names []string
	err = db.Select(&names, "SELECT name FROM place LIMIT 10")

	//PREPARED statement
	stmt, err := db.Prepare(`SELECT * FROM place WHERE telecode = $1`)
	row = stmt.QueryRow(65)
	fmt.Println("prepared row : ", row)

}
