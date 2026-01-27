package handlers

import (
	"fmt"
	"net/http"
)

func SeatsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Get Seats request")
	case http.MethodPost:
		fmt.Fprintln(w, "Seats POST Request")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}
