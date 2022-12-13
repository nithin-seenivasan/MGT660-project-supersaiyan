package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// BVK - This function ideally should be under controllers.go
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
	// BVK - move checking business logic to separate functions in event_models.go

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
			errorMessage = "Unable to accept input. Please check the entered data. Note: only png|jpg|jpeg|gif|gifv images are supported"
			var redirectURL string = "/events/new"
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		//Insert Kim Kardashian as default attendee
		// BVK - Move setting up default events into another function and call it from here
		rsvpData := Rsvp{
			EventID:      newID,
			EmailAddress: "kim.kardashian@yale.edu",
		}
		_, rsvpError := addRSVP(rsvpData)

		rsvpData2 := Rsvp{
			EventID:      newID,
			EmailAddress: "kyle.jensen@yale.edu",
		}
		_, rsvpError2 := addRSVP(rsvpData2)

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

	//Redirect back to events/new
	var redirectURL string = "/events/new"
	http.Redirect(w, r, redirectURL, http.StatusFound)

	// BVK - remove commented code below

	//Display the create page with the concatenated error Message (containing aggregate of all error messages)
	//tmpl["create-error"].Execute(w, errorMessage)

}
