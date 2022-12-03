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

	event_id, err := strconv.Atoi(path.Base(r.URL.Path))
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad email address", http.StatusBadRequest)
			return
		}

		email_address := r.FormValue("email")

		rsvpData := Rsvp{
			Event_ID:      event_id,
			Email_address: email_address,
		}

		confirmation_code, database_err := addRSVP(rsvpData)

		if database_err != nil {
			//Error here comes from the INSERT SQL statement - display the following message
			errors := "This is a Yale exclusive event. Please enter a @yale.edu email address only"
			ContextData := setupEventContextData(w, event_id, "", errors)
			tmpl["events"].Execute(w, ContextData)
		} else {
			ContextData := setupEventContextData(w, event_id, confirmation_code[:7], "")
			tmpl["events"].Execute(w, ContextData)
		}

	} else {
		contextData := setupEventContextData(w, event_id, "", "")
		tmpl["events"].Execute(w, contextData)
	}

}

func createController(w http.ResponseWriter, r *http.Request) {
	if errorMessage != "" {
		//Display the create page with the concatenated error Message (containing aggregate of all error messages)
		tmpl["create-error"].Execute(w, errorMessage)
		errorMessage = ""
	} else {
		tmpl["create"].Execute(w, "")
	}

}

func setupEventContextData(w http.ResponseWriter, event_id int, confirmation_code string, errors string) EventContextData {

	Requested_Event, err := getEventByID(event_id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return EventContextData{}
	}

	RSVP_List, _ := getRSVPByID(event_id)

	contextData := EventContextData{
		Event:             Requested_Event,
		Rsvp_data:         RSVP_List,
		Confirmation_Code: confirmation_code,
		Errors:            errors,
	}

	return contextData
}

func donateController(w http.ResponseWriter, r *http.Request) {

	tmpl["donate"].Execute(w, "nothing")
}
