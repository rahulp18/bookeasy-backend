package repository

import (
	"context"
	"database/sql"

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
