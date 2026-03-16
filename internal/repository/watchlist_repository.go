package repository

import (
	"context"
	"fmt"
	"stock_screener/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchlistRepository struct {
	db *pgxpool.Pool
}

func NewWatchlistRepository(db *pgxpool.Pool) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}

// Add saves a stock to the user's watchlist.
// Safe to call multiple times — uses ON CONFLICT DO NOTHING.
func (r *WatchlistRepository) Add(ctx context.Context, userID, stockID uuid.UUID) (*domain.WatchlistEntry, error) {
	query := `
		INSERT INTO watchlists (user_id, stock_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, stock_id) DO NOTHING
		RETURNING id, user_id, stock_id, created_at
	`

	var entry domain.WatchlistEntry
	err := r.db.QueryRow(ctx, query, userID, stockID).Scan(
		&entry.ID, &entry.UserID, &entry.StockID, &entry.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("watchlist repository: Add: %w", err)
	}

	return &entry, nil
}

// Remove deletes a stock from the user's watchlist.
func (r *WatchlistRepository) Remove(ctx context.Context, userID, stockID uuid.UUID) error {
	query := `DELETE FROM watchlists WHERE user_id = $1 AND stock_id = $2`

	if _, err := r.db.Exec(ctx, query, userID, stockID); err != nil {
		return fmt.Errorf("watchlist repository: Remove: %w", err)
	}

	return nil
}

// GetByUser returns all watchlist items for a user, enriched with stock details.
func (r *WatchlistRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.WatchlistItem, error) {
	query := `
		SELECT
			wl.id, wl.user_id, wl.stock_id, wl.created_at,
			s.id, s.ticker, s.company_name, s.sector, s.industry,
			s.market_cap, s.pe_ratio, s.eps, s.revenue_growth,
			s.ranking_score, s.created_at, s.updated_at,
			COALESCE(p.price, 0),
			COALESCE(p.change_percent, 0),
			COALESCE(p.updated_at, NOW())
		FROM watchlists wl
		JOIN stocks s       ON s.id = wl.stock_id
		LEFT JOIN stock_prices p ON p.stock_id = s.id
		WHERE wl.user_id = $1
		ORDER BY wl.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("watchlist repository: GetByUser: %w", err)
	}
	defer rows.Close()

	var items []domain.WatchlistItem
	for rows.Next() {
		var item domain.WatchlistItem
		var sd domain.StockDetail
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.StockID, &item.CreatedAt,
			&sd.ID, &sd.Ticker, &sd.CompanyName, &sd.Sector, &sd.Industry,
			&sd.MarketCap, &sd.PERatio, &sd.EPS, &sd.RevenueGrowth,
			&sd.RankingScore, &sd.CreatedAt, &sd.UpdatedAt,
			&sd.Price, &sd.ChangePercent, &sd.PriceUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("watchlist repository: GetByUser scan: %w", err)
		}
		item.Stock = sd
		items = append(items, item)
	}

	return items, rows.Err()
}
