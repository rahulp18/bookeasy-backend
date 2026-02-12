package models

import "time"

type Event struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes"`
	CreatedAt       time.Time `json:"created_at"`
}

type EventDetails struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"duration_minutes"`
	Shows           []Show
	CreatedAt       time.Time `json:"created_at"`
}
type EventUpdateRequest struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	DurationMinutes *string `json:"duration_minutes"`
}
