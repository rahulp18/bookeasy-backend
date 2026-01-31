package services

import (
	"context"
	"errors"

	"github.com/rahulp18/bookeasy-backend/internal/models"
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

func (bs *BookingService) ConfirmBooking(ctx context.Context, userID, bookingID string) error {
	//   validate
	ok, err := bs.bookingRepo.IsBookingPendingAndOwner(ctx, userID, bookingID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Booking not found or already proceed")

	}
	return bs.bookingRepo.ConfirmBooking(ctx, bookingID)
}

func (bs *BookingService) GetBookingDetails(ctx context.Context, bookingID, userID string) (models.BookingDetails, error) {
	return bs.bookingRepo.GetBookingDetails(ctx, bookingID, userID)
}
func (bs *BookingService) GetAllBookings(ctx context.Context, userID string) ([]models.UserBookingSummary, error) {
	return bs.bookingRepo.GetUserBookingSummery(ctx, userID)
}
