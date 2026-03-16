package service

import (
	"context"
	"errors"
	"fmt"
	"stock_screener/internal/domain"
	"stock_screener/internal/repository"

	"github.com/google/uuid"
)

// ErrStockNotFound is returned when a ticker doesn't exist in the DB.
var ErrStockNotFound = errors.New("stock not found")

// ErrNoStocksAvailable is returned when the discover pool is exhausted.
var ErrNoStocksAvailable = errors.New("no stocks available to discover")

type StockService struct {
	stockRepo *repository.StockRepository
}

func NewStockService(stockRepo *repository.StockRepository) *StockService {
	return &StockService{stockRepo: stockRepo}
}

// Discover returns the next stock a user should see.
// Excludes stocks the user has already swiped or watchlisted.
// Returns ErrNoStocksAvailable if the user has seen everything.
func (s *StockService) Discover(ctx context.Context, userID uuid.UUID) (*domain.StockDetail, error) {
	stocks, err := s.stockRepo.Discover(ctx, userID, 1)
	if err != nil {
		return nil, fmt.Errorf("stock service: Discover: %w", err)
	}

	if len(stocks) == 0 {
		return nil, ErrNoStocksAvailable
	}

	return &stocks[0], nil
}

// GetByTicker returns a stock's full detail by ticker symbol.
func (s *StockService) GetByTicker(ctx context.Context, ticker string) (*domain.StockDetail, error) {
	stock, err := s.stockRepo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("stock service: GetByTicker: %w", ErrStockNotFound)
	}

	return stock, nil
}

// GetHistory returns OHLCV history for charting.
func (s *StockService) GetHistory(ctx context.Context, ticker string, rangeStr string) ([]domain.StockHistory, error) {
	stock, err := s.stockRepo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, ErrStockNotFound
	}

	history, err := s.stockRepo.GetHistory(ctx, stock.ID, rangeStr)
	if err != nil {
		return nil, fmt.Errorf("stock service: GetHistory: %w", err)
	}

	return history, nil
}

// GetByTickers returns multiple stocks by ticker symbols.
func (s *StockService) GetByTickers(ctx context.Context, tickers []string) ([]domain.StockDetail, error) {
	stocks, err := s.stockRepo.GetByTickers(ctx, tickers)
	if err != nil {
		return nil, fmt.Errorf("stock service: GetByTickers: %w", err)
	}

	return stocks, nil
}

// Search stocks by ticker or company name.
func (s *StockService) Search(ctx context.Context, query string) ([]domain.StockDetail, error) {
	stocks, err := s.stockRepo.Search(ctx, query, 20)
	if err != nil {
		return nil, fmt.Errorf("stock service: Search: %w", err)
	}

	return stocks, nil
}
