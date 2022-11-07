package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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
		fmt.Println("Error 404 page does not exist")
		os.Exit(1)
	}

	type EventContextData struct {
		Event     Event
		Rsvp_data []string
	}

	Requested_Event, err := getEventByID(fetch_id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	contextData := EventContextData{
		Event:     Requested_Event,
		Rsvp_data: []string{"dummyemail1@yale.edu", "dummyemail2@gmail.com"},
	}

	tmpl["events"].Execute(w, contextData)
}

func createController(w http.ResponseWriter, r *http.Request) {
	tmpl["create"].Execute(w, "")
}
