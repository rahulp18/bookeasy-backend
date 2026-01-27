package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := repository.GetAllUsers()

		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Context-Type", "application.json")
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		fmt.Fprintln(w, "POST Request")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}
