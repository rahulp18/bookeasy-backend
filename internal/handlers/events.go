package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/services"
)

type AdminEventHandler struct {
	eventService *services.AdminEventService
}

func NewAdminEventHandler(eventService *services.AdminEventService) *AdminEventHandler {
	return &AdminEventHandler{eventService: eventService}
}

type CreateEventRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"duration_minutes"`
}

func (h *AdminEventHandler) EventActionHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *AdminEventHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := h.eventService.GetAllEvents(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	case http.MethodPost:
		var input CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		id, err := h.eventService.CreateEvent(r.Context(), models.Event{
			Title:           input.Title,
			Description:     input.Description,
			DurationMinutes: input.DurationMinutes,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"event_id": id,
		})

	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)

	}

}
