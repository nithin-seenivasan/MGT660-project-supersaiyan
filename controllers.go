package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"
)

// Home Page
func indexController(w http.ResponseWriter, r *http.Request) {

	type indexContextData struct {
		Events []Event
		Today  time.Time
	}

	theEvents, err := getAllEvents()
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	contextData := indexContextData{
		Events: theEvents,
		Today:  time.Now(),
	}

	tmpl["index"].Execute(w, contextData)
}

// About Page
func aboutController(w http.ResponseWriter, r *http.Request) {

	tmpl["about"].Execute(w, "nothing")
}

// Show each event as a page - RSVP logic in here
func eventsController(w http.ResponseWriter, r *http.Request) {

	eventID, err := strconv.Atoi(path.Base(r.URL.Path))
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	contextData := setupEventContextData(w, eventID, "", "")

	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad email address", http.StatusBadRequest)
			return
		}

		emailAddress := r.FormValue("email")

		rsvpData := Rsvp{
			EventID:      eventID,
			EmailAddress: emailAddress,
		}

		confirmationCode, databaseErr := addRSVP(rsvpData)
		errors := ""
		if confirmationCode != "" {
			confirmationCode = confirmationCode[:7]
		}

		if databaseErr != nil {
			//Error here comes from the INSERT SQL statement - display the following message
			errors = "This is a Yale exclusive event. Please enter a @yale.edu email address only"
		}

		contextData = setupEventContextData(w, eventID, confirmationCode, errors)

	}

	tmpl["events"].Execute(w, contextData)
}

// Handles the GET and POST routes for /events/new
func addNewEventController(w http.ResponseWriter, r *http.Request) {

	//GET route - just show
	if r.Method == http.MethodGet {
		tmpl["create"].Execute(w, "")
		return
	}

	// The submit button links to /events/new as a POST request, which is routed by controller in routes.go to come HERE.
	// Extract the form's POST variables here. Write to DB. Then IF all is OK, redirect to confirmation page
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Extracts values from the POST request
	title := r.FormValue("title")
	location := r.FormValue("location")
	image := r.FormValue("image")
	date := r.FormValue("date")

	//Check business logic in a separate function in event_models.go
	errorMessage, parsedDate := checkEventData(title, location, image, date)

	//Create a Event element with the new variables
	newEvent := Event{
		Title:    title,
		Location: location,
		Image:    image,
		Date:     parsedDate,
	}

	//If no Error, execute this -
	if errorMessage == "" {
		//Insert into Postgres DB.
		newID, err := addEvent(newEvent)
		if err != nil {
			//Error here comes from the INSERT SQL statement - display the following message
			errorMessage = "Unable to accept input. Please check the entered data. Note: only png|jpg|jpeg|gif|gifv images are supported"
			var redirectURL string = "/events/new"
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		//Insert Kim Kardashian and Kyle Jensen as default attendees
		//Uses an external function
		var rsvpData, rsvpData2 Rsvp
		rsvpData, rsvpData2 = defaultRSVPEmails(newID)

		//Add RSVP for default attendees
		_, rsvpError := addRSVP(rsvpData)
		_, rsvpError2 := addRSVP(rsvpData2)

		//Handle Errors
		if rsvpError != nil || rsvpError2 != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		//Redirect to the events/{ID} page, on successful entry into DB
		var redirectURL string = "/events/" + strconv.Itoa(newID)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}
	//Else if Error -
	errorMessage = errorMessage + " Please try again!"

	//Show error message on page
	tmpl["create"].Execute(w, errorMessage)

}

// Returns the thank you for donating page
func donateController(w http.ResponseWriter, r *http.Request) {
	tmpl["donate"].Execute(w, "nothing")
}
