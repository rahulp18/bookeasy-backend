package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rahulp18/bookeasy-backend/internal/services"
)

type SeatHandler struct {
	seatService *services.SeatService
}

func NewSeatHandler(seatService *services.SeatService) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
	}
}
func (sh *SeatHandler) SeatIdHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	fmt.Print(parts)
	showID := parts[2]
	switch r.Method {
	case http.MethodGet:
		seats, err := sh.seatService.GetShowSeats(r.Context(), showID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(seats)
	case http.MethodPost:
		fmt.Fprintln(w, "Seats POST Request")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}
