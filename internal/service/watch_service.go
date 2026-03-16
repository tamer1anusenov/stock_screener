package service

import (
	"context"
	"fmt"
	"stock_screener/internal/domain"
	"stock_screener/internal/repository"

	"github.com/google/uuid"
)

type WatchlistService struct {
	watchlistRepo *repository.WatchlistRepository
}

func NewWatchlistService(watchlistRepo *repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{watchlistRepo: watchlistRepo}
}

// GetWatchlist returns all saved stocks for a user, enriched with stock details.
func (s *WatchlistService) GetWatchlist(ctx context.Context, userID uuid.UUID) ([]domain.WatchlistItem, error) {
	items, err := s.watchlistRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("watchlist service: GetWatchlist: %w", err)
	}

	return items, nil
}

// Remove deletes a stock from the user's watchlist.
func (s *WatchlistService) Remove(ctx context.Context, userID, stockID uuid.UUID) error {
	if err := s.watchlistRepo.Remove(ctx, userID, stockID); err != nil {
		return fmt.Errorf("watchlist service: Remove: %w", err)
	}

	return nil
}
