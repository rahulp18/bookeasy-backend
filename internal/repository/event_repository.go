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
