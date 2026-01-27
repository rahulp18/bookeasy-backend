package models

type Seat struct {
	ID     string `json:"id"`
	ShowID string `json:"show_id"`
	Number string `json:"number"`
	Status string `json:"status"`
}
