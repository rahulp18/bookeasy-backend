package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/rahulp18/bookeasy-backend/internal/models"
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
func (br *BookingRepository) IsBookingPendingAndOwner(ctx context.Context, userID, bookingID string) (bool, error) {
	var exists bool
	err := br.db.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE id=$1 AND user_id=$2 AND status='pending'
		)`,
		bookingID,
		userID,
	).Scan(&exists)
	return exists, err
}

func (br *BookingRepository) GetBookingDetails(ctx context.Context, bookingID, userID string) (models.BookingDetails, error) {
	var booking models.BookingDetails

	// 1) Fetch booking + show + event (single row)
	err := br.db.QueryRowContext(ctx, `
		SELECT
			b.id,
			b.status,
			b.created_at,

			s.id,
			s.venue,
			s.start_time,
			s.end_time,

			e.id,
			e.title,
			e.description,
			e.duration_minutes
		FROM bookings b
		JOIN shows s ON b.show_id = s.id
		JOIN events e ON s.event_id = e.id
		WHERE b.id = $1 AND b.user_id = $2
	`, bookingID, userID).Scan(
		&booking.BookingID,
		&booking.Status,
		&booking.CreatedAt,

		&booking.Show.ShowID,
		&booking.Show.Venue,
		&booking.Show.StartTime,
		&booking.Show.EndTime,

		&booking.Show.Event.EventID,
		&booking.Show.Event.Title,
		&booking.Show.Event.Description,
		&booking.Show.Event.DurationMinutes,
	)

	if err != nil {
		return models.BookingDetails{}, err
	}

	// 2) Fetch seats (multiple rows)
	rows, err := br.db.QueryContext(ctx, `
		SELECT
			ss.id AS show_seat_id,
			se.id AS seat_id,
			se.seat_row,
			se.seat_number
		FROM booking_seats bs
		JOIN show_seats ss ON bs.show_seat_id = ss.id
		JOIN seats se ON ss.seat_id = se.id
		WHERE bs.booking_id = $1
		ORDER BY se.seat_row, se.seat_number
	`, bookingID)

	if err != nil {
		return models.BookingDetails{}, err
	}
	defer rows.Close()

	seats := []models.SeatInfo{}

	for rows.Next() {
		var seat models.SeatInfo
		var seatRow string
		var seatNumber string

		err := rows.Scan(
			&seat.ShowSeatID,
			&seat.SeatID,
			&seatRow,
			&seatNumber,
		)
		if err != nil {
			return models.BookingDetails{}, err
		}

		seat.RowLabel = seatRow
		seat.SeatNo = seatNumber

		seats = append(seats, seat)
	}

	if err := rows.Err(); err != nil {
		return models.BookingDetails{}, err
	}

	booking.Seats = seats
	return booking, nil
}

func (br *BookingRepository) GetUserBookingSummery(ctx context.Context, userID string) ([]models.UserBookingSummary, error) {

	query := `
		SELECT
			b.id AS booking_id,
			b.status,
			b.created_at,
			e.title AS event_title,
			s.venue,
			s.start_time
		FROM bookings b
		JOIN shows s ON b.show_id = s.id
		JOIN events e ON s.event_id = e.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`
	rows, err := br.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookings []models.UserBookingSummary
	for rows.Next() {
		var b models.UserBookingSummary
		err := rows.Scan(
			&b.BookingID,
			&b.Status,
			&b.CreatedAt,
			&b.EventTitle,
			&b.Venue,
			&b.StartTime,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, err
}
