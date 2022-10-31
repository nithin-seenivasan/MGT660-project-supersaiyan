package main

import (
	"net/http"
)

func aboutController(w http.ResponseWriter, r *http.Request) {

	//  type indexContextData struct {
	// 	Events []Event
	// 	Today  time.Time
	// }

	// theEvents, err := getAllEvents()
	// if err != nil {
	// 	http.Error(w, "database error", http.StatusInternalServerError)
	// 	return
	// }

	//contextData := 0
	// 	Events: theEvents,
	// 	Today:  time.Now(),
	// }

	tmpl["about"].Execute(w, "nothing")
}
