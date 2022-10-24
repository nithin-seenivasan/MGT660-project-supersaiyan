package main

import (
	"log"
	"os"
	"time"

	_ "embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// Event - encapsulates information about an event
type Event struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Location  string    `json:"location"`
	Image     string    `json:"image"`
	Date      time.Time `json:"date"`
	Attending []string  `json:"attending"`
}

// getEventByID - returns the event that has the
// specified id and an error if there was a database err.
func getEventByID(id int) (Event, error) {
	var e Event
	query := `SELECT id,title,location,image,date FROM events WHERE id=$1`
	row := db.QueryRow(query, id)
	err := row.Scan(&e.ID, &e.Title, &e.Location, &e.Image, &e.Date)
	return e, err
}

// getAllEvents - returns a slice of all events and an
// error in case of database error.
func getAllEvents() ([]Event, error) {
	var events []Event
	query := `SELECT id,title,location,image,date FROM events`
	rows, err := db.Query(query)
	if err != nil {
		return events, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Location, &e.Image, &e.Date)
		if err != nil {
			return events, err
		}
		events = append(events, e)
	}
	return events, nil
}

func addEvent(event Event) (int, error) {
	insertStatement := `
		INSERT INTO events (title, location, image, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	// Use the global db object
	newID := 0
	err := db.QueryRow(insertStatement, event.Title, event.Location, event.Image, event.Date).Scan(&newID)
	return newID, err
}

//go:embed init-schema.sql
var makeSchema string

// The above comment causes the `init-schema.sql` file to be
// embedded in the binary. See https://golang.org/pkg/embed/
// The makeSchema variable is a string that contains the
// contents of the `init-schema.sql` file.

// init is run once when this file is first loaded. See
// https://golang.org/doc/effective_go.html#init
// https://medium.com/golangspec/init-functions-in-go-eac191b3860a
func init() {
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		// e.g.: YourUserName:YourPassword@YourHost:5432/YourDatabase
		panic("You must supply the DATABASE_URL")
	}
	var dbErr error

	// I suggest you use sqlx instead of golang's built-in sql
	// library. sqlx does not add much, but what it does add
	// is useful. See https://jmoiron.github.io/sqlx/
	db, dbErr = sqlx.Open("postgres", databaseURL)
	if dbErr != nil {
		panic("Could not connect to database")
	}

	_, err := db.Exec(makeSchema)
	if err != nil {
		log.Panicln("Could not create database schema:", err)
	}
}
