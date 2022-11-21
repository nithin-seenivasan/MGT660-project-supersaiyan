package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var errorMessage string //global variable
func addNewEventController(w http.ResponseWriter, r *http.Request) {
	// The submit button links to /events/new-event-created, which is rendered by routes.go to come HERE.
	// Extract the form's POST variables here. (see the YT video). Write to DB. Then IF all is OK, do this below
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

	//Create a Error Message based on conditions

	if len(title) < 6 {
		errorMessage = "Title is invalid!"
	}

	if len(location) < 6 {
		errorMessage = errorMessage + " Location is invalid!"
	}

	//Use DB "CHECK" statement for image URL

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
			tmpl["create-error"].Execute(w, "Unable to accept input. Please check the entered data. Note: only png|jpg|jpeg|gif|gifv images are supported")
			return
		}
		//Insert Kim Kardashian as default attendee
		rsvp_data := Rsvp{
			Event_ID:      newID,
			Email_address: "kim.kardashian@yale.edu",
		}
		_, rsvp_error := addRSVP(rsvp_data)
		if rsvp_error != nil {
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

	//Redirect back to events/new
	var redirectURL string = "/events/new"
	http.Redirect(w, r, redirectURL, http.StatusFound)

	//Display the create page with the concatenated error Message (containing aggregate of all error messages)
	//tmpl["create-error"].Execute(w, errorMessage)

}
