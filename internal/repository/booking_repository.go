package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}
func (br *BookingRepository) CreateBooking(ctx context.Context, userID string, showID string, showSeatIDs []string) (string, error) {
	if len(showSeatIDs) == 0 {
		return "", errors.New("No seats are selected")
	}
	tx, err := br.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	// 1. Lock seat
	query := `SELECT id FROM show_seats
	        WHERE id = ANY($1) AND status ='available'
			FOR UPDATE
	`
	rows, err := tx.QueryContext(ctx, query, pq.Array(showSeatIDs))
	if err != nil {
		return "", err
	}

	defer rows.Close()

	locked := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return "", err
		}
		locked = append(locked, id)
	}

	if len(locked) != len(showSeatIDs) {
		return "", errors.New("one or more seats are not available")
	}
	_, err = tx.ExecContext(ctx, `UPDATE show_seats SET status ='locked',locked_at =NOW() WHERE id =ANY($1)`,
		pq.Array(showSeatIDs),
	)
	if err != nil {
		return "", err
	}
	// 2 CREATE BOOKING
	var bookingID string
	err = tx.QueryRowContext(ctx, `INSERT INTO bookings(user_id,show_id,status)
	 VALUES ($1,$2,'pending')
	 RETURNING id
	`, userID, showID).Scan(&bookingID)
	if err != nil {
		return "", err
	}
	// 4 Insert booking_seats
	for _, seatId := range showSeatIDs {
		_, err := tx.ExecContext(ctx, `INSERT INTO booking_seats(booking_id,show_seat_id) VALUES ($1,$2)`, bookingID, seatId)
		if err != nil {
			return "", err
		}
	}
	// 5 COMMIT
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return bookingID, nil
}

func (br *BookingRepository) ConfirmBooking(ctx context.Context, bookingID string) error {
	tx, err := br.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//    1. CONFIRM BOOKING
	_, err = tx.ExecContext(ctx,
		`UPDATE bookings
		 SET status='confirmed'
		 WHERE id=$1 AND status='pending'`,
		bookingID,
	)
	if err != nil {
		return err
	}
	// 2 MARKED SEAT BOOKED
	_, err = tx.ExecContext(ctx,
		`UPDATE show_seats
		 SET status='booked', locked_at=NULL
		 WHERE id IN (
			SELECT show_seat_id
			FROM booking_seats
			WHERE booking_id=$1
		 )`,
		bookingID,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (br *BookingRepository) CancelBooking(ctx context.Context, bookingID string) error {
	tx, err := br.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//   CANCEL BOOKING
	_, err = tx.ExecContext(ctx, `UPDATE bookings SET status='cancelled' WHERE id=$1`, bookingID)
	if err != nil {
		return err
	}
	// 2 RELEASE SEATS
	_, err = tx.ExecContext(ctx, `UPDATE show_seats SET status='available',locked_at=NULL WHERE id IN(
	 SELECT show_seat_id
	 FROM booking_seats
     WHERE booking_ID=$1
	)`, bookingID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (br *BookingRepository) ReleaseExpiredSeats(ctx context.Context) error {
	tx, err := br.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `UPDATE show_seats SET status='available', locked_at=NULL WHERE ...`)
	if err != nil {
		return err
	}
	return tx.Commit()
}
