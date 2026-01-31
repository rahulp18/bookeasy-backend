package models

import "time"

type UserBookingSummary struct {
	BookingID string    `json:"booking_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	EventTitle string    `json:"event_title"`
	Venue      string    `json:"venue"`
	StartTime  time.Time `json:"start_time"`
}
