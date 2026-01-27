package routes

import (
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.HandleFunc("/events", handlers.EventsHandler)
	mux.HandleFunc("/seats", handlers.SeatsHandler)
	mux.HandleFunc("/shows", handlers.ShowsHandler)
}
