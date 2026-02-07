package models

import "time"

type BookingDetails struct {
	BookingID string    `json:"booking_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	Show  ShowDetails `json:"show"`
	Seats []SeatInfo  `json:"seats"`
}

type ShowDetails struct {
	ShowID    string    `json:"show_id"`
	Venue     string    `json:"venue"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	Event EventInfo `json:"event"`
}

type EventInfo struct {
	EventID         string `json:"event_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"duration_minutes"`
}

type SeatInfo struct {
	ShowSeatID string `json:"show_seat_id"`
	SeatID     string `json:"seat_id"`
	SeatNo     string `json:"seat_no"`
	RowLabel   string `json:"row_label"`
}
