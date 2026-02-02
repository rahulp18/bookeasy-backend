package services

import (
	"context"
	"errors"

	"github.com/rahulp18/bookeasy-backend/internal/models"
	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

type SeatService struct {
	seatRepo *repository.SeatRepository
}

func NewSeatService(seatRepo *repository.SeatRepository) *SeatService {
	return &SeatService{
		seatRepo: seatRepo,
	}
}
func (ss *SeatService) GetShowSeats(ctx context.Context, showID string) ([]models.ShowSeatsResponse, error) {
	if showID == "" {
		return nil, errors.New("show_id is required")
	}
	return ss.seatRepo.GetSeatsByShowID(ctx, showID)
}
