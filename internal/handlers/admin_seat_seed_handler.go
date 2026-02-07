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
func (h *AdminSeatSeedHandler) SeedShowSeats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	showID := parts[3]
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
}
