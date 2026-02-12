package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rahulp18/bookeasy-backend/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}
func (er *EventRepository) CreateEvent(ctx context.Context, e models.Event) (string, error) {
	var id string
	query := `
	INSERT INTO events(title,description,duration_minutes)
	VALUES ($1,$2,$3)
	RETURNING id
	`
	err := er.db.QueryRowContext(ctx, query, e.Title, e.Description, e.DurationMinutes).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (er *EventRepository) FetchAllEvents(ctx context.Context) ([]models.Event, error) {
	query := `SELECT id,title,description,duration_minutes,created_at FROM events ORDER BY created_at DESC `
	rows, err := er.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	events := make([]models.Event, 0)
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.DurationMinutes,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (er *EventRepository) GetEventDetails(ctx context.Context, eventID string) (models.EventDetails, error) {
	var details models.EventDetails
	eventQuery := `SELECT id,title,description,duration_minutes,created_at FROM events
	        WHERE id=$1       
	`
	err := er.db.QueryRowContext(ctx, eventQuery, eventID).Scan(&details.ID, &details.Title, &details.Description, &details.DurationMinutes, &details.CreatedAt)
	if err != nil {
		return models.EventDetails{}, err
	}
	showsQuery := `SELECT id,event_id,venue,start_time,end_time,created_at FROM shows WHERE event_id=$1 ORDER BY start_time ASC`
	rows, err := er.db.QueryContext(ctx, showsQuery, eventID)
	if err != nil {
		return models.EventDetails{}, err
	}
	defer rows.Close()
	shows := make([]models.Show, 0)
	for rows.Next() {
		var s models.Show
		err := rows.Scan(
			&s.ID,
			&s.EventId,
			&s.Venue,
			&s.StartTime,
			&s.EndTime,
			&s.CreatedAt,
		)
		if err != nil {
			return models.EventDetails{}, err
		}
		shows = append(shows, s)
	}
	if err := rows.Err(); err != nil {
		return models.EventDetails{}, err
	}
	details.Shows = shows
	return details, nil
}

func (er *EventRepository) DeleteEvent(ctx context.Context, eventID string) error {
	tx, err := er.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//  1 Delete shows first

	_, err = tx.ExecContext(ctx, `DELETE FROM shows WHERE event_id = $1`, eventID)
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM events WHERE id=$1`, eventID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}
func (er *EventRepository) UpdateEvent(ctx context.Context, eventID string, eventData models.EventUpdateRequest) error {
	fmt.Println(eventData)
	query := `
		UPDATE events
		SET 
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			duration_minutes = COALESCE($3, duration_minutes)
		WHERE id = $4
	`
	res, err := er.db.ExecContext(ctx, query, eventData.Title, eventData.Description, eventData.DurationMinutes, eventID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
