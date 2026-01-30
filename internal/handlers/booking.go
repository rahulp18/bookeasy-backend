package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/middleware"
	"github.com/rahulp18/bookeasy-backend/internal/services"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

type CreateBookingSeatsRequest struct {
	ShowID      string   `json:"show_id"`
	ShowSeatIDs []string `json:"show_seat_ids"`
}
type CreateBookingResponse struct {
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
}

func (bh *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req CreateBookingSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "show_id and show_seat_ids are required", http.StatusBadRequest)
		return
	}
	bookingID, err := bh.bookingService.CreateBooking(
		r.Context(),
		userID,
		req.ShowID,
		req.ShowSeatIDs,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateBookingResponse{
		BookingID: bookingID,
		Status:    "pending",
	})

}
