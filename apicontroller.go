package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

func apiEventListController(w http.ResponseWriter, r *http.Request) {
	type apiResponse struct {
		Events []Event `json:"events"`
	}
	allEvents, _ := getAllEvents()
	a := apiResponse{Events: allEvents}
	json, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func apiEventDetailController(w http.ResponseWriter, r *http.Request) {

	eventID, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "Bad event ID", http.StatusBadRequest)
		return
	}
	event, wasFound := getEventByID(eventID)
	if wasFound != nil {
		http.Error(w, "No such Event", http.StatusNotFound)
		return
	}
	json, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
