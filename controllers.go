package main

import (
	"net/http"
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
	tmpl["events"].Execute(w, "") //sending empty string, form uses it for potential error message.
}

func createController(w http.ResponseWriter, r *http.Request) {
	tmpl["create"].Execute(w, "nothing")
}

func addNewEventToDb(w http.ResponseWriter, r *http.Request) {

	// The submit button links to /events/new-event-created, which is rendered by routes.go to come HERE.
	// Extract the form's POST variables here. Then IF all is OK, do this below
	tmpl["post-creation"].Execute(w, "nothing")

	//If not, send back form with injected error message!
}
