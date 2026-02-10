package services

import (
	"context"
	"errors"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

type AdminEventService struct {
	eventRepo *repository.EventRepository
}

func NewAdminEventService(eventRepo *repository.EventRepository) *AdminEventService {
	return &AdminEventService{eventRepo: eventRepo}
}

func (aes *AdminEventService) CreateEvent(ctx context.Context, e models.Event) (string, error) {
	if e.Title == "" {
		return "", errors.New("Title is required")
	}
	if e.DurationMinutes <= 0 {
		return "", errors.New("duration_minutes must be >0")
	}
	return aes.eventRepo.CreateEvent(ctx, e)
}

func (aes *AdminEventService) GetAllEvents(ctx context.Context) ([]models.Event, error) {

	return aes.eventRepo.FetchAllEvents(ctx)
}

func (aes *AdminEventService) GetEventDetails(ctx context.Context, eventID string) (models.EventDetails, error) {
	return aes.eventRepo.GetEventDetails(ctx, eventID)
}
func (aes *AdminEventService) DeleteEvent(ctx context.Context, eventID string) error {
	return aes.eventRepo.DeleteEvent(ctx, eventID)
}
func (aes *AdminEventService) UpdateEvent(ctx context.Context, eventID string, eventData models.EventUpdateRequest) error {
	return aes.eventRepo.UpdateEvent(ctx, eventID, eventData)
}
