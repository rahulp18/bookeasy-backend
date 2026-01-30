package repository

import (
	"context"
	"log"
	"time"
)

func StartSeatExpiryWorker(ctx context.Context, repo *BookingRepository) {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := repo.ReleaseExpiredSeats(ctx); err != nil {
					log.Println("Seat expiry error", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
