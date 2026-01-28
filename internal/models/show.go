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
