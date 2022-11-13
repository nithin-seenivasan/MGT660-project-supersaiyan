package main

import (
	"testing"
	"time"
)

// Test to check if the database accepts an invalid Title
func TestAddEventTitle(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")
	if err != nil {
		println("Time parsing error")
	}

	var testEvents = []Event{
		{
			Title:    "FOUR",
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

	testEvent := Event{
		Title:    "Party In New Haven",
		Location: "N",
		Image:    "https://i.imgur.com/l3aFizL.jpeg",
		Date:     testDate,
	}

	got, err := addEvent(testEvent)

	if err != nil {
		t.Errorf("Invalid Location - Not accepted in Database. Official error message: %q", err)
	} else {
		println(got)
	}
}

// Test to check if the database accepts an invalid image
func TestAddEventImage(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")
	if err != nil {
		println("Time parsing error")
	}

	testEvent := Event{
		Title:    "Party In New Haven",
		Location: "N",
		Image:    "1",
		Date:     testDate,
	}

	got, err := addEvent(testEvent)

	if err != nil {
		t.Errorf("Invalid Image format - Not accepted in Database. Official error message: %q", err)
	} else {
		println(got)
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
