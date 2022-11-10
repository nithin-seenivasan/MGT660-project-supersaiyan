package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func apiEventListController(w http.ResponseWriter, r *http.Request) {
	type apiResponse struct {
		Events []Event `json:"events"`
	}
	allEvents, err := getAllEvents()
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
	eventID, err := strconv.Atoi(chi.URLParam(r, "eventID"))
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
