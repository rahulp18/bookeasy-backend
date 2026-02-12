package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/services"
)

type AdminSeatSeedHandler struct {
	seatSeedService *services.AdminSeatSeedService
}

func NewAdminSeatSeedHandler(seatSeedService *services.AdminSeatSeedService) *AdminSeatSeedHandler {
	return &AdminSeatSeedHandler{
		seatSeedService: seatSeedService,
	}
}

func (h *AdminSeatSeedHandler) HandleShowSeatsRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	showID := parts[3]
	switch r.Method {
	case http.MethodGet:
		showDetails, err := h.seatSeedService.GetShowDetails(r.Context(), showID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(showDetails)
	case http.MethodPost:
		var input models.SeedSeatRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		err := h.seatSeedService.SeedShowSeats(r.Context(), showID, input.Rows, input.SeatsPerRow)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Seats seeded successfully",
		})
	case http.MethodDelete:
		err := h.seatSeedService.DeleteShow(r.Context(), showID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Deleted",
		})
	case http.MethodPatch:
		var input models.ShowUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.seatSeedService.UpdateShow(r.Context(), showID, input)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Updated",
		})
	default:
		http.Error(w, "Method not allowed Rahul", http.StatusMethodNotAllowed)
	}

}
