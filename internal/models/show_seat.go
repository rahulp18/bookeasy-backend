package models

import "time"

type ShowSeat struct {
	ID       string    `json:"id"`
	ShowID   string    `json:"show_id"`
	SeatID   string    `json:"seat_id"`
	Status   string    `json:"status"`
	LockedAt time.Time `json:"locked_at"`
}
