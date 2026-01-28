package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/middleware"
	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

type ProfileResponse struct {
	Success bool         `json:"success"`
	Data    *models.User `json:"data"`
}

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

func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application.json")
	if http.MethodGet != r.Method {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	user_id, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || user_id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserById(user_id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Context-Type", "application.json")
	json.NewEncoder(w).Encode(ProfileResponse{
		Success: true,
		Data:    user,
	})
}
