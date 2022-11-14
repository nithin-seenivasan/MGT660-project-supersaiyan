package main

import (
	"github.com/go-chi/chi"
)

func createRoutes() chi.Router {
	// We're using chi as the router. You'll want to read
	// the documentation https://github.com/go-chi/chi
	// so that you can capture parameters like /events/5
	// or /api/events/4 -- where you want to get the
	// event id (5 and 4, respectively).

	r := chi.NewRouter()
	r.Get("/", indexController)
	r.Get("/about", aboutController)
	r.Get("/events/rsvp/{event_id:[0-9]+}", addrsvpController)
	r.Get("/events/{event_id:[0-9]+}", eventsController)
	r.Get("/events/new", createController)
	r.Get("/events/new-event-created", addNewEventController)
	addStaticFileServer(r, "/static/", "staticfiles")
	return r
}
