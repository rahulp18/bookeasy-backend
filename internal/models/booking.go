package models

import "time"

type Booking struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ShowID    string    `json:"show_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
