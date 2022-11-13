package main

import (
	"testing"
	"time"
)

//Run the tests by running "go test" in CMD

// Test to check if the database accepts an invalid Title
func TestAddEventTitle(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")
	if err != nil {
		println("Time parsing error")
	}

	var testEvents = []Event{
		{
			Title:    "FOUR", //4 Characters
			Location: "New Haven",
			Image:    "https://i.imgur.com/l3aFizL.jpeg",
			Date:     testDate},
		{
			Title:    "52characters52characters52characters52characters5252", //52 characters
			Location: "New Haven",
			Image:    "https://i.imgur.com/l3aFizL.jpeg",
			Date:     testDate},
	}

	for _, event := range testEvents {
		if got, err := addEvent(event); err != nil {
			t.Errorf("Invalid Title - Not accepted in Database. Official error message: %q ||| Input Title: %q", err, event.Title)
		} else {
			println(got)
		}
	}

}

// Test to check if the database accepts an invalid location
func TestAddEventLocation(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")
	if err != nil {
		println("Time parsing error")
	}

	var testEvents = []Event{
		{
			Title:    "Test Event from Unit Testing",
			Location: "FOUR", // 4 characters
			Image:    "https://i.imgur.com/l3aFizL.jpeg",
			Date:     testDate},
		{
			Title:    "Test Event from Unit Testing",
			Location: "52characters52characters52characters52characters5252", //52 characters
			Image:    "https://i.imgur.com/l3aFizL.jpeg",
			Date:     testDate},
	}

	for _, event := range testEvents {
		if got, err := addEvent(event); err != nil {
			t.Errorf("Invalid Location - Not accepted in Database. Official error message: %q ||| Input Title: %q", err, event.Title)
		} else {
			println(got)
		}
	}
}

// Test to check if the database accepts an invalid image
func TestAddEventImage(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")
	if err != nil {
		println("Time parsing error")
	}

	var testEvents = []Event{
		{
			Title:    "Test Event from Unit Testing",
			Location: "New Haven",
			Image:    "https://i.imgur.com/l3aFizL.pdf", //PDF
			Date:     testDate},
		{
			Title:    "Test Event from Unit Testing",
			Location: "New Haven",
			Image:    "http", //4 characters
			Date:     testDate},
		{
			Title:    "Test Event from Unit Testing",
			Location: "New Haven",
			Image:    "52characters52characters52characters52characters5252.jpg", //52 characters
			Date:     testDate},
		{
			Title:    "Test Event from Unit Testing",
			Location: "New Haven",
			Image:    "htt://i.imgur.com/l3aFizL.jpeg", //No HTTP or HTTPS
			Date:     testDate},
		{
			Title:    "Test Event from Unit Testing",
			Location: "New Haven",
			Image:    "https://i.imgur.com/l3aFizL", //No file extension
			Date:     testDate},
	}

	for _, event := range testEvents {
		if got, err := addEvent(event); err != nil {
			t.Errorf("Invalid Image Link - Not accepted in Database. Official error message: %q ||| Input Title: %q", err, event.Title)
		} else {
			println(got)
		}
	}
}

// Test to check if the database accepts an valid Event
func TestAddEventValid(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")
	if err != nil {
		println("Time parsing error")
	}

	testEvent := Event{
		Title:    "Party In New Haven",
		Location: "New Haven",
		Image:    "https://i.imgur.com/l3aFizL.jpeg",
		Date:     testDate,
	}

	got, err := addEvent(testEvent)

	if err != nil {
		t.Errorf("Invalid Image format - Not accepted in Database. Official error message: %q", err)
	} else {
		println("TestAddEventValid test successful. Could write to Database - The ID number of the row added:", got)
	}
}
