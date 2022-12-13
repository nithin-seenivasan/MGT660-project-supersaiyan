package main

import (
	"errors"
	"net/http"
	"path"
	"strconv"
	"time"
)

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

func aboutController(w http.ResponseWriter, r *http.Request) {

	tmpl["about"].Execute(w, "nothing")
}

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

// BVK - this function shouldn't exist
func createController(w http.ResponseWriter, r *http.Request) {
	if errorMessage != "" {
		//Display the create page with the concatenated error Message (containing aggregate of all error messages)
		tmpl["create"].Execute(w, errorMessage) //Changed to send it to same Create page - Bala's elegant solution used here - Error only shows IF errorMessage != empty
		errorMessage = ""
	} else {
		tmpl["create"].Execute(w, "")
	}

}

func donateController(w http.ResponseWriter, r *http.Request) {

	tmpl["donate"].Execute(w, "nothing")
}
