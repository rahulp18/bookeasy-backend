package services

import (
	"context"
	"errors"
	"time"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

type AdminShowService struct {
	showRepo *repository.ShowRepository
}

func NewAdminShowService(showRepo *repository.ShowRepository) *AdminShowService {
	return &AdminShowService{showRepo: showRepo}
}

func (ass *AdminShowService) CreateShow(ctx context.Context, s models.Show) (string, error) {
	if s.EventId == "" {
		return "", errors.New("event_id is required")
	}
	if s.Venue == "" {
		return "", errors.New("Venue is required")
	}
	if s.StartTime.IsZero() || s.EndTime.IsZero() {
		return "", errors.New("start_time and end_time are required")
	}
	if s.EndTime.Before(s.StartTime) {
		return "", errors.New("end_time connot be before start_time")
	}
	if s.StartTime.Before(time.Now().Add(-1 * time.Hour)) {
		return "", errors.New("start_time looks invalid (past)")
	}
	return ass.showRepo.CreateShow(ctx, s)
}
