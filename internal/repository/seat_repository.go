package repository

import (
	"context"
	"database/sql"

	"github.com/rahulp18/bookeasy-backend/internal/models"
)

type SeatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) *SeatRepository {
	return &SeatRepository{
		db: db,
	}
}

func (sr *SeatRepository) GetSeatsByShowID(ctx context.Context, showID string) ([]models.ShowSeatsResponse, error) {
	query := `
		SELECT
			ss.id,
			s.seat_row,
			s.seat_number,
			ss.status
		FROM show_seats ss
		JOIN seats s ON ss.seat_id = s.id
		WHERE ss.show_id = $1
		ORDER BY s.seat_row, s.seat_number
	`
	rows, err := sr.db.QueryContext(ctx, query, showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.ShowSeatsResponse

	for rows.Next() {
		var seat models.ShowSeatsResponse
		err := rows.Scan(
			&seat.ShowSeatID,
			&seat.RowLevel,
			&seat.SeatNumber,
			&seat.Status,
		)
		if err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}
