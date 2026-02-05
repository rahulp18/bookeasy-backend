package repository

import (
	"context"
	"database/sql"

	"github.com/rahulp18/bookeasy-backend/internal/models"
)

type ShowRepository struct {
	db *sql.DB
}

func NewShowRepository(db *sql.DB) *ShowRepository {
	return &ShowRepository{
		db: db,
	}
}

func (sr *ShowRepository) CreateShow(ctx context.Context, s models.Show) (string, error) {
	var id string
	query := `
	INSERT INTO shows(event_id,venue,start_time,end_time)
	VALUES ($1,$2,$3,$4)
	RETURNING id
	`
	err := sr.db.QueryRowContext(ctx, query, s.EventId, s.Venue, s.StartTime, s.EndTime).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}
