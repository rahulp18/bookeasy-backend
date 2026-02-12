package services

import (
	"context"
	"errors"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

type AdminSeatSeedService struct {
	seedRepo *repository.SeatSeedRepository
}

func NewAdminSeatSeedService(seedRepo *repository.SeatSeedRepository) *AdminSeatSeedService {
	return &AdminSeatSeedService{seedRepo: seedRepo}
}

func (as *AdminSeatSeedService) SeedShowSeats(ctx context.Context, showID string, rows []string, seatsPerRow int) error {
	if showID == "" {
		return errors.New("show_id is required")
	}
	if len(rows) == 0 {
		return errors.New("rows is required")
	}
	if seatsPerRow <= 0 {
		return errors.New("seats_per_row must be > 0")
	}

	return as.seedRepo.SeedShowSeats(ctx, showID, rows, seatsPerRow)
}

func (as *AdminSeatSeedService) GetShowDetails(ctx context.Context, showID string) (models.ShowDetailsRes, error) {
	return as.seedRepo.GetShowDetails(ctx, showID)
}
func (as *AdminSeatSeedService) DeleteShow(ctx context.Context, showID string) error {
	return as.seedRepo.DeleteShow(ctx, showID)
}
func (as *AdminSeatSeedService) UpdateShow(ctx context.Context, showID string, showInput models.ShowUpdateRequest) error {
	return as.seedRepo.UpdateShow(ctx, showID, showInput)
}
