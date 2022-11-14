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

func getEventData(w http.ResponseWriter, event_id int) EventContextData {

	Requested_Event, err := getEventByID(event_id)
	if err != nil {
		http.Error(w, "Event not present", http.StatusInternalServerError)
		return EventContextData{}
	}

	RSVP_List, err := getRSVPByID(event_id)

	contextData := EventContextData{
		Event:     Requested_Event,
		Rsvp_data: RSVP_List,
	}

	return contextData
}

func eventsController(w http.ResponseWriter, r *http.Request) {

	fetch_id, err := strconv.Atoi(path.Base(r.URL.Path))
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tmpl["events"].Execute(w, getEventData(w, fetch_id))
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

	EventData := getEventData(w, event_id)

	database_err := addRSVP(rsvpData)
	if database_err != nil {
		//Error here comes from the INSERT SQL statement - display the following message
		tmpl["rsvp_error"].Execute(w, EventData)
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

	type RSVPContextData struct {
		Event     Event
		RSVP_data []string
		Code      string
	}

	contextData := RSVPContextData{
		Event:     EventData.Event,
		RSVP_data: EventData.Rsvp_data,
		Code:      Code,
	}
	tmpl["rsvp"].Execute(w, contextData)
}
