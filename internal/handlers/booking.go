package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rahulp18/bookeasy-backend/internal/middleware"
	"github.com/rahulp18/bookeasy-backend/internal/models"
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
type BookingDetailsResponse struct {
	Success bool                  `json:"success"`
	Data    models.BookingDetails `json:"data"`
}
type BookingsResponse struct {
	Success bool                        `json:"success"`
	Data    []models.UserBookingSummary `json:"data"`
}

func (bh *BookingHandler) HandleBookings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPost:
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
	case http.MethodGet:
		data, err := bh.bookingService.GetAllBookings(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BookingsResponse{
			Success: true,
			Data:    data,
		})
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}

}

func (bh *BookingHandler) BookingActions(w http.ResponseWriter, r *http.Request) {
	// GET BOOKING ID FROM PARAMS
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	bookingID := parts[2]
	action := ""
	if len(parts) > 3 {
		action = parts[3]
	}
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		if action == "confirm" {
			err := bh.bookingService.ConfirmBooking(r.Context(), userID, bookingID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "booking confirmed",
			})
		} else {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		}
	case http.MethodGet:
		data, err := bh.bookingService.GetBookingDetails(r.Context(), bookingID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(BookingDetailsResponse{
			Success: true,
			Data:    data,
		})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
