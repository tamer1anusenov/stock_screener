package service

import (
	"context"
	"errors"
	"fmt"

	"stock_screener/internal/domain"
	"stock_screener/internal/repository"

	"github.com/google/uuid"
)

// ErrAlreadySwiped is returned when a user tries to swipe a stock they've already swiped.
var ErrAlreadySwiped = errors.New("already swiped this stock")

type SwipeService struct {
	swipeRepo     *repository.SwipeRepository
	watchlistRepo *repository.WatchlistRepository
}

func NewSwipeService(
	swipeRepo *repository.SwipeRepository,
	watchlistRepo *repository.WatchlistRepository,
) *SwipeService {
	return &SwipeService{
		swipeRepo:     swipeRepo,
		watchlistRepo: watchlistRepo,
	}
}

// Swipe records a swipe action.
// If direction is "right", the stock is automatically added to the watchlist.
func (s *SwipeService) Swipe(ctx context.Context, userID, stockID uuid.UUID, direction domain.SwipeDirection) (*domain.Swipe, error) {
	// Guard: prevent duplicate swipes
	already, err := s.swipeRepo.HasSwiped(ctx, userID, stockID)
	if err != nil {
		return nil, fmt.Errorf("swipe service: %w", err)
	}
	if already {
		return nil, ErrAlreadySwiped
	}

	// Record the swipe
	swipe, err := s.swipeRepo.Create(ctx, userID, stockID, direction)
	if err != nil {
		return nil, fmt.Errorf("swipe service: %w", err)
	}

	// Auto-add to watchlist on right swipe
	if direction == domain.SwipeRight {
		if _, err := s.watchlistRepo.Add(ctx, userID, stockID); err != nil {
			// Non-fatal: log it but don't fail the swipe
			fmt.Printf("swipe service: failed to add to watchlist: %v\n", err)
		}
	}

	return swipe, nil
}
