package main

import (
	"testing"
	"time"
)

func TestAddEventTitle(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")

	testEvent := Event{
		Title:    "A",
		Location: "New Haven",
		Image:    "https://i.imgur.com/l3aFizL.jpeg",
		Date:     testDate,
	}

	got, err := addEvent(testEvent)

	if err != nil {
		t.Errorf("Invalid Title - Not accepted in Database. Official error message: %q", err)
	} else {
		println(got)
	}
}

func TestAddEventLocation(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")

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

func TestAddEventImage(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")

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

func TestAddEventValid(t *testing.T) {

	testDate, err := time.Parse("2006-01-02T15:04", "2006-01-02T15:04")

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
