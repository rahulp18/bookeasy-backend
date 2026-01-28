package models

type Seat struct {
	ID         string `json:"id"`
	SeatRow    string `json:"seat_row"`
	SeatNumber string `json:"seat_number"`
}
