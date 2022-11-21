package main

import (
	"log"
	"net/http"
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

type Rsvp struct {
	Event_ID      int    `json:"eventid"`
	Email_address string `json:"email"`
}

type EventContextData struct {
	Event             Event
	Rsvp_data         []string
	Confirmation_Code string
	Errors            string
}

// getEventByID - returns the event that has the
// specified id and an error if there was a database err.
func getEventByID(id int, w http.ResponseWriter) (Event, error) {
	var e Event
	query := `SELECT id,title,location,image,date FROM events WHERE id=$1`
	row := db.QueryRow(query, id)
	err := row.Scan(&e.ID, &e.Title, &e.Location, &e.Image, &e.Date)

	//Added this because the API was not sending RSVP data - just showing attending: "null"
	//for each event ID in events, search rsvp for all attending and add it to events.attanding
	//Just made use of Bala's function, and linked it to Mike's already created API functions
	//This bit of code acts like a bridge between those two
	var contextData EventContextData
	contextData = setupEventContextData(w, e.ID, "", "")
	e.Attending = contextData.Rsvp_data

	return e, err
}

// getAllEvents - returns a slice of all events and an
// error in case of database error.
func getAllEvents(w http.ResponseWriter) ([]Event, error) {
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

		//Added this because the API was not sending RSVP data - just showing attending: "null"
		//for each event ID in events, search rsvp for all attending and add it to events.attanding
		//Just made use of Bala's function, and linked it to Mike's already created API functions
		//This bit of code acts like a bridge between those two
		var contextData EventContextData
		contextData = setupEventContextData(w, e.ID, "", "")
		e.Attending = contextData.Rsvp_data

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
	newID := 0
	err := db.QueryRow(insertStatement, event.Title, event.Location, event.Image, event.Date).Scan(&newID)
	return newID, err
}

func addRSVP(rsvp_data Rsvp) (string, error) {
	insertStatement := `
		INSERT INTO rsvp (event_id, email_address)		
		VALUES ($1, $2)
		RETURNING confirmation_code;
	`
	code := "00000"
	err := db.QueryRow(insertStatement, rsvp_data.Event_ID, rsvp_data.Email_address).Scan(&code)
	return code, err
}

func getRSVPByID(id int) ([]string, error) {
	rsvp_list := []string{}
	query := `SELECT email_address FROM rsvp WHERE event_id=$1`
	rows, err := db.Query(query, id)
	if err != nil {
		return rsvp_list, err
	}
	defer rows.Close()
	for rows.Next() {
		var e string
		err := rows.Scan(&e)
		if err != nil {
			return rsvp_list, err
		}
		rsvp_list = append(rsvp_list, e)
	}
	return rsvp_list, nil
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
