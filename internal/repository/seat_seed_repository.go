package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type SeatSeedRepository struct {
	db *sql.DB
}

func NewSeatSeedRepository(db *sql.DB) *SeatSeedRepository {
	return &SeatSeedRepository{db: db}
}

func (ssr *SeatSeedRepository) SeedShowSeats(ctx context.Context, showID string, rows []string, seatsPerRow int) error {
	tx, err := ssr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//   1. CREATE SEATS IF NOT EXISTS
	// 2. CREATE SHOW_SEATS FOR EACH SEAT

	for _, row := range rows {
		for i := 1; i <= seatsPerRow; i++ {
			var seatID string

			err := tx.QueryRowContext(ctx, `
		INSERT INTO seats(seat_row,seat_number)
		VALUES ($1,$2)
		ON CONFLICT(seat_row,seat_number)
		DO UPDATE SET seat_row = EXCLUDED.seat_row
        RETURNING id
		`, row, fmt.Sprintf("%d", i)).Scan(&seatID)

			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, `
	INSERT INTO show_seats(show_id,seat_id,status)
	VALUES ($1,$2,'available')
    ON CONFLICT(show_id, seat_id) DO NOTHING
	`, showID, seatID)
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
