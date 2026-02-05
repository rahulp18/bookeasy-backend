package services

import (
	"context"
	"errors"

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
