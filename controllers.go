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

	fetch_id, err := strconv.Atoi(path.Base(r.URL.Path))
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	type EventContextData struct {
		Event     Event
		Rsvp_data []string
	}

	Requested_Event, err := getEventByID(fetch_id)
	if err != nil {
		http.Error(w, "database error 1", http.StatusInternalServerError)
		return
	}

	RSVP_List, err := getRSVPByID(fetch_id)
	// if err != nil {
	// 	http.Error(w, "database error", http.StatusInternalServerError)
	// 	return
	// }

	contextData := EventContextData{
		Event:     Requested_Event,
		Rsvp_data: RSVP_List,
	}

	tmpl["events"].Execute(w, contextData)
}

func createController(w http.ResponseWriter, r *http.Request) {
	tmpl["create"].Execute(w, "")
}

func addrsvpController(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad email address", http.StatusBadRequest)
		return
	}
	email_address := r.FormValue("email")

	event_id, err := strconv.Atoi(path.Base(r.URL.Path))
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, "Bad event ID", http.StatusBadRequest)
		return
	}

	rsvpData := Rsvp{
		Event_ID:      event_id,
		Email_address: email_address,
	}

	database_err := addRSVP(rsvpData)
	if database_err != nil {
		//Error here comes from the INSERT SQL statement - display the following message
		tmpl["events"].Execute(w, "This is a Yale exclusive event. Please enter a @yale.edu email address only")
		return
	}

	var redirectURL string = "/events/" + strconv.Itoa(event_id)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
