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

// Rsvp - encapsulates RSVP information - email address and eventID
type Rsvp struct {
	EventID      int    `json:"eventid"`
	EmailAddress string `json:"email"`
}

// Function that returns default attendees to every new event created
func defaultRSVPEmails(newID int) (Rsvp, Rsvp) {
	rsvpKim := Rsvp{
		EventID:      newID,
		EmailAddress: "kim.kardashian@yale.edu",
	}

	rsvpKyle := Rsvp{
		EventID:      newID,
		EmailAddress: "kyle.jensen@yale.edu",
	}

	return rsvpKim, rsvpKyle
}

// Business logic - checks for various conditions and returns an error message
func checkEventData(title string, location string, image string, date string) (errorMessage string, parsedDate time.Time) {

	//Check lengths of title and location
	if len(title) < 6 || len(title) > 49 {
		errorMessage = "Title is invalid!"
	}

	if len(location) < 6 || len(location) > 49 {
		errorMessage = errorMessage + " Location is invalid!"
	}

	//Parses date string to time.Time element
	parsedDate, err := time.Parse("2006-01-02T15:04", date)
	if err != nil {
		errorMessage = errorMessage + " Date is invalid!"
	}

	//Compare dates
	today := time.Now()
	dateComparison := parsedDate.After(today)

	if !dateComparison {
		errorMessage = errorMessage + " Date is in the past!"
	}

	return errorMessage, parsedDate
}

// EventContextData - encapsulates Event Context information - Event, RSVP, ConfirmationCode and Errors
type EventContextData struct {
	Event            Event
	RsvpData         []string
	ConfirmationCode string
	Errors           string
}

// getEventByID - returns the event that has the
// specified id and an error if there was a database err.
func getEventByID(id int) (Event, error) {
	var e Event
	query := `SELECT id,title,location,image,date FROM events WHERE id=$1`
	row := db.QueryRow(query, id)
	err := row.Scan(&e.ID, &e.Title, &e.Location, &e.Image, &e.Date)
	//e.Attending = []string{"ABC", "DEF"}

	RSVPList, _ := getRSVPByID(e.ID)
	e.Attending = RSVPList

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

		//Added this because the API was not sending RSVP data - just showing attending: "null"
		//for each event ID in events, search rsvp for all attending and add it to events.attanding
		//Just made use of Bala's function, and linked it to Mike's already created API functions
		//This bit of code acts like a bridge between those two
		RSVPList, _ := getRSVPByID(e.ID)
		e.Attending = RSVPList

		events = append(events, e)
	}

	return events, nil
}

// addEvent - Store event information in database, returns an error
// if event cannot be inserted
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

// addRSVP - store RSVP data containing eventID and emailAddress
// returns confirmationCode and error
func addRSVP(RsvpData Rsvp) (string, error) {
	insertStatement := `
		INSERT INTO rsvp (event_id, email_address)		
		VALUES ($1, $2)
		RETURNING confirmation_code;
	`
	code := ""
	err := db.QueryRow(insertStatement, RsvpData.EventID, RsvpData.EmailAddress).Scan(&code)
	return code, err
}

// getRSVPbyID - fetch RSVP for a particular eventID
// returns list of RSVPs and errors
func getRSVPByID(id int) ([]string, error) {
	rsvpList := []string{}
	query := `SELECT email_address FROM rsvp WHERE event_id=$1`
	rows, err := db.Query(query, id)
	if err != nil {
		return rsvpList, err
	}
	defer rows.Close()
	for rows.Next() {
		var e string
		err := rows.Scan(&e)
		if err != nil {
			return rsvpList, err
		}
		rsvpList = append(rsvpList, e)
	}
	return rsvpList, nil
}

func setupEventContextData(w http.ResponseWriter, eventID int, confirmationCode string, errors string) EventContextData {

	requestedEvent, err := getEventByID(eventID)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return EventContextData{}
	}

	RSVPList, _ := getRSVPByID(eventID)

	contextData := EventContextData{
		Event:            requestedEvent,
		RsvpData:         RSVPList,
		ConfirmationCode: confirmationCode,
		Errors:           errors,
	}

	return contextData
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
