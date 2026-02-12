package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rahulp18/bookeasy-backend/internal/models"
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
func (sr *SeatSeedRepository) UpdateShow(ctx context.Context, showID string, showInput models.ShowUpdateRequest) error {
	//  1st find the Show from Show ID
	var show models.Show
	err := sr.db.QueryRowContext(ctx, `SELECT id,venue,start_time,end_time FROM shows WHERE id=$1`, showID).Scan(&show.ID, &show.Venue, &show.StartTime, &show.EndTime)

	if err != nil {
		return err
	}
	if show.ID != showID {
		return errors.New("Show not found with this ID")
	}
	//  Now Let's Update Show
	res, err := sr.db.ExecContext(ctx, `UPDATE shows SET venue=COALESCE($1,venue),start_time=COALESCE($2,start_time),end_time=COALESCE($3,end_time) WHERE id=$4`, showInput.Venue, showInput.StartTime, showInput.EndTime, showID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Failed to update")
	}
	return nil
}

func (sr *SeatSeedRepository) DeleteShow(ctx context.Context, showID string) error {
	tx, err := sr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//  DELETE SHOW SEATS FIRST
	_, err = tx.ExecContext(ctx, `DELETE FROM show_seats WHERE show_id=$1 `, showID)
	if err != nil {
		return err
	}
	// DELETE SHOW NOW
	res, err := tx.ExecContext(ctx, `DELETE FROM shows WHERE is=$1`, showID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Failed to delete show")
	}
	return tx.Commit()
}
func (sr *SeatSeedRepository) GetShowDetails(
	ctx context.Context,
	showID string,
) (models.ShowDetailsRes, error) {

	var showDetails models.ShowDetailsRes

	// 1) Fetch show basic info
	query := `
		SELECT id, venue, start_time, end_time, created_at
		FROM shows
		WHERE id = $1
	`

	err := sr.db.QueryRowContext(ctx, query, showID).Scan(
		&showDetails.ID,
		&showDetails.Venue,
		&showDetails.StartTime,
		&showDetails.EndTime,
		&showDetails.CreatedAt,
	)
	if err != nil {
		return models.ShowDetailsRes{}, err
	}

	// 2) Fetch show seats
	seatQuery := `
		SELECT
			ss.id,
			ss.status,
			s.id,
			s.seat_row,
			s.seat_number
		FROM show_seats ss
		JOIN seats s ON ss.seat_id = s.id
		WHERE ss.show_id = $1
		ORDER BY s.seat_row, s.seat_number
	`

	rows, err := sr.db.QueryContext(ctx, seatQuery, showID)
	if err != nil {
		return models.ShowDetailsRes{}, err
	}
	defer rows.Close()

	showSeats := []models.ShowSeatInfo{}

	for rows.Next() {
		var seat models.ShowSeatInfo

		err := rows.Scan(
			&seat.ID,
			&seat.Status,
			&seat.SeatID,
			&seat.SeatRow,
			&seat.SeatNumber,
		)
		if err != nil {
			return models.ShowDetailsRes{}, err
		}

		showSeats = append(showSeats, seat)
	}

	if err := rows.Err(); err != nil {
		return models.ShowDetailsRes{}, err
	}

	showDetails.ShowSeats = showSeats

	return showDetails, nil
}
