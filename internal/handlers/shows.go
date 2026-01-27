package handlers

import (
	"fmt"
	"net/http"
)

func ShowsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Get shows request")
	case http.MethodPost:
		fmt.Fprintln(w, "shows POST Request")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}
