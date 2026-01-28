package models

type ShowSeat struct {
	ID     string `json:"id"`
	ShowID string `json:"show_id"`
	SeatID string `json:"seat_id"`
	Status string `json:"status"`
}
