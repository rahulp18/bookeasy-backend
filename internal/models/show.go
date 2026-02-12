package models

import "time"

type Show struct {
	ID        string    `json:"id"`
	EventId   string    `json:"event_id"`
	Venue     string    `json:"venue"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
}

type ShowUpdateRequest struct {
	Venue     *string `json:"venue"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
}

type ShowDetailsRes struct {
	ID        string         `json:"id"`
	Venue     string         `json:"venue"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	CreatedAt time.Time      `json:"created_at"`
	ShowSeats []ShowSeatInfo `json:"show_seats"`
}

type ShowSeatInfo struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	SeatID     string `json:"seat_id"`
	SeatRow    string `json:"seat_row"`
	SeatNumber string `json:"seat_number"`
}
