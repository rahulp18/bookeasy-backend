package services

import (
	"context"
	"errors"

	"github.com/rahulp18/bookeasy-backend/internal/repository"
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
	}
}

func (bs *BookingService) CreateBooking(
	ctx context.Context,
	userId string,
	showID string,
	showSeatIDs []string,
) (string, error) {
	//    Business validation
	if len(showSeatIDs) == 0 {
		return "", errors.New("no seats selected")
	}
	if len(showSeatIDs) > 6 {
		return "", errors.New("Cannot book more then 6 seats")
	}
	bookingID, err := bs.bookingRepo.CreateBooking(ctx, userId, showID, showSeatIDs)
	if err != nil {
		return "", err
	}
	return bookingID, nil
}
