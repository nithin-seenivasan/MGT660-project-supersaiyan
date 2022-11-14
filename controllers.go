package main

import (
	"crypto/rand"
	"encoding/hex"
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

func getEventData(w http.ResponseWriter, event_id int) (Event, []string) {

	Requested_Event, err := getEventByID(event_id)
	if err != nil {
		http.Error(w, "Event not present", http.StatusInternalServerError)
		return Event{}, []string{}
	}

	RSVP_List, _ := getRSVPByID(event_id)

	return Requested_Event, RSVP_List
}

func eventsController(w http.ResponseWriter, r *http.Request) {

	fetch_id, err := strconv.Atoi(path.Base(r.URL.Path))
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	Requested_Event, RSVP_List := getEventData(w, fetch_id)

	type ContextData struct {
		Event     Event
		Rsvp_data []string
	}

	EventContextData := ContextData{
		Event:     Requested_Event,
		Rsvp_data: RSVP_List,
	}

	tmpl["events"].Execute(w, EventContextData)
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

	Requested_Event, RSVP_List := getEventData(w, event_id)

	type ContextData struct {
		Event     Event
		Rsvp_data []string
		Code      string
	}

	EventContextData := ContextData{
		Event:     Requested_Event,
		Rsvp_data: RSVP_List,
		Code:      "",
	}

	database_err := addRSVP(rsvpData)
	if database_err != nil {
		//Error here comes from the INSERT SQL statement - display the following message
		tmpl["rsvp_error"].Execute(w, EventContextData)
		return
	}

	// var redirectURL string = "/events/" + strconv.Itoa(event_id)
	// http.Redirect(w, r, redirectURL, http.StatusFound)
	ByteCode := make([]byte, 6)
	_, error := rand.Read(ByteCode)
	Code := ""
	if error != nil {
		Code = "code not found"
	} else {
		Code = hex.EncodeToString(ByteCode)
	}

	RSVP_ContextData := ContextData{
		Event:     Requested_Event,
		Rsvp_data: RSVP_List,
		Code:      Code,
	}
	tmpl["rsvp"].Execute(w, RSVP_ContextData)
}
