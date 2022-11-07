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
	tmpl["create"].Execute(w, "")
}

func addNewEventToDb(w http.ResponseWriter, r *http.Request) {

	// The submit button links to /events/new-event-created, which is rendered by routes.go to come HERE.
	// Extract the form's POST variables here. (see the YT video). Write to DB. Then IF all is OK, do this below
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Extracts values from the POST request
	title := r.FormValue("title")
	location := r.FormValue("location")
	image := r.FormValue("image")
	date := r.FormValue("date")

	//Create a Error Message based on conditions
	var errorMessage string
	if len(title) < 2 {
		errorMessage = "Title is invalid!"
	}

	if len(location) < 2 {
		errorMessage = errorMessage + " Location is invalid!"
	}

	//Use DB "CHECK" statement for image URL

	//Parses date string to time.Time element
	parsedDate, err := time.Parse("2006-01-02T15:04", date)
	if err != nil {
		errorMessage = errorMessage + " Date is invalid!"
	}

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
			tmpl["create"].Execute(w, "Unable to accept input. Please check the entered data. Note: only png|jpg|jpeg|gif|gifv images are supported")
			return
		}
		tmpl["post-creation"].Execute(w, newID)
		return
	}
	errorMessage = errorMessage + " Please try again!"

	//Display the create page with the concatenated error Message (containing aggregate of all error messages)
	tmpl["create"].Execute(w, errorMessage)

}
