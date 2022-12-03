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
	r.Post("/events/{event_id:[0-9]+}", eventsController)
	r.Get("/events/{event_id:[0-9]+}", eventsController)
	r.Get("/events/new", createController)
	r.Post("/events/new-event-created", addNewEventController) //Form returns a post action, so the Chi controller needs POST to respond to it
	r.Get("/api/events", apiEventListController)
	r.Get("/api/events/{event_id}", apiEventDetailController)
	addStaticFileServer(r, "/static/", "staticfiles")
	r.Get("/events/{event_id:[0-9]+/donate}", donateController)
	return r
}
