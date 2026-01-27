package handlers

import (
	"fmt"
	"net/http"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Get events request")
	case http.MethodPost:
		fmt.Fprintln(w, "Events POST Request")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}
