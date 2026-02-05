package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/services"
)

type AdminShowHandler struct {
	showService *services.AdminShowService
}
type CreateShowRequest struct {
	EventID   string `json:"event_id"`
	Venue     string `json:"venue"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func NewAdminShowHandler(showService *services.AdminShowService) *AdminShowHandler {
	return &AdminShowHandler{showService: showService}
}
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

func (h *AdminShowHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	var input CreateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	startTime, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		http.Error(w, "Invalid start_time format (use RFC3339)", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, input.EndTime)
	if err != nil {
		http.Error(w, "Invalid end_time format (use RFC3339)", http.StatusBadRequest)
		return
	}
	id, err := h.showService.CreateShow(r.Context(), models.Show{
		EventId:   input.EventID,
		Venue:     input.Venue,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"show_id": id,
	})
}
