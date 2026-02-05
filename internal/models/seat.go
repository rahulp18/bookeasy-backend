package models

type Seat struct {
	ID         string `json:"id"`
	SeatRow    string `json:"seat_row"`
	SeatNumber string `json:"seat_number"`
}

type ShowSeatsResponse struct {
	ShowSeatID string `json:"show_seat_id"`
	RowLevel   string `json:"row_level"`
	SeatNumber string `json:"seat_number"`
	Status     string `json:"status"`
}

type SeedSeatRequest struct {
	Rows        []string `json:"rows"`
	SeatsPerRow int      `json:"seats_per_show"`
}
