package models

import "time"

type Show struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	EventId   string    `json:"event_id"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	CreatedAt time.Time `json:"created_at"`
}
