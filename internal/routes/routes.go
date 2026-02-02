package routes

import (
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/db"
	"github.com/rahulp18/bookeasy-backend/internal/handlers"
	"github.com/rahulp18/bookeasy-backend/internal/middleware"
	"github.com/rahulp18/bookeasy-backend/internal/repository"
	"github.com/rahulp18/bookeasy-backend/internal/services"
)

func Register(mux *http.ServeMux) {
	bookingRepository := repository.NewBookingRepository(db.DB)
	bookingService := services.NewBookingService(bookingRepository)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	seatRepository := repository.NewSeatRepository(db.DB)
	seatService := services.NewSeatService(seatRepository)
	seatHandler := handlers.NewSeatHandler(seatService)

	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.Handle("/profile", middleware.Auth(http.HandlerFunc(handlers.Profile)))
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/events", handlers.EventsHandler)
	// mux.HandleFunc("/seats", handlers.SeatsHandler)
	mux.HandleFunc("/shows/", seatHandler.SeatIdHandler)
	// BOOKING ROUTES
	mux.Handle("/bookings", middleware.Auth(http.HandlerFunc(bookingHandler.HandleBookings)))
	mux.Handle("/bookings/", middleware.Auth(http.HandlerFunc(bookingHandler.BookingActions)))
}
