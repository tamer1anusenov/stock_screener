package repository

import (
	"context"
	"fmt"
	"stock_screener/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StockRepository struct {
	db *pgxpool.Pool
}

func NewStockRepository(db *pgxpool.Pool) *StockRepository {
	return &StockRepository{db: db}
}

// GetByTicker returns a stock's fundamentals by ticker symbol.
func (r *StockRepository) GetByTicker(ctx context.Context, ticker string) (*domain.StockDetail, error) {
	query := `
		SELECT
			s.id, s.ticker, s.company_name, s.sector, s.industry,
			s.market_cap, s.pe_ratio, s.eps, s.revenue_growth,
			s.ranking_score, s.created_at, s.updated_at,
			COALESCE(p.price, 0),
			COALESCE(p.change_percent, 0),
			COALESCE(p.updated_at, NOW())
		FROM stocks s
		LEFT JOIN stock_prices p ON p.stock_id = s.id
		WHERE s.ticker = $1
	`

	var sd domain.StockDetail
	err := r.db.QueryRow(ctx, query, ticker).Scan(
		&sd.ID, &sd.Ticker, &sd.CompanyName, &sd.Sector, &sd.Industry,
		&sd.MarketCap, &sd.PERatio, &sd.EPS, &sd.RevenueGrowth,
		&sd.RankingScore, &sd.CreatedAt, &sd.UpdatedAt,
		&sd.Price, &sd.ChangePercent, &sd.PriceUpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("stock repository: GetByTicker %s: %w", ticker, err)
	}
	return &sd, nil
}

// Discover returns stocks the user has not yet swiped or watchlisted,
// ordered by ranking_score descending.
func (r *StockRepository) Discover(ctx context.Context, userID uuid.UUID, limit int) ([]domain.StockDetail, error) {
	query := `
		SELECT
			s.id, s.ticker, s.company_name, s.sector, s.industry,
			s.market_cap, s.pe_ratio, s.eps, s.revenue_growth,
			s.ranking_score, s.created_at, s.updated_at,
			COALESCE(p.price, 0),
			COALESCE(p.change_percent, 0),
			COALESCE(p.updated_at, NOW())
		FROM stocks s
		LEFT JOIN stock_prices p  ON p.stock_id = s.id
		LEFT JOIN swipes sw       ON sw.stock_id = s.id AND sw.user_id = $1
		LEFT JOIN watchlists wl   ON wl.stock_id = s.id AND wl.user_id = $1
		WHERE sw.stock_id IS NULL
		  AND wl.stock_id IS NULL
		ORDER BY s.ranking_score DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("stock repository: Discover: %w", err)
	}
	defer rows.Close()

	var stocks []domain.StockDetail
	for rows.Next() {
		var sd domain.StockDetail
		if err := rows.Scan(
			&sd.ID, &sd.Ticker, &sd.CompanyName, &sd.Sector, &sd.Industry,
			&sd.MarketCap, &sd.PERatio, &sd.EPS, &sd.RevenueGrowth,
			&sd.RankingScore, &sd.CreatedAt, &sd.UpdatedAt,
			&sd.Price, &sd.ChangePercent, &sd.PriceUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("stock repository: Discover scan: %w", err)
		}
		stocks = append(stocks, sd)
	}

	return stocks, rows.Err()
}

// GetHistory returns OHLCV candles for a stock within the given range.
// rangeStr accepted values: "1m", "3m", "6m", "1y", "5y"
func (r *StockRepository) GetHistory(ctx context.Context, stockID uuid.UUID, rangeStr string) ([]domain.StockHistory, error) {
	interval := rangeToInterval(rangeStr)

	query := `
		SELECT stock_id, date, open, high, low, close, volume
		FROM stock_history
		WHERE stock_id = $1
		  AND date >= NOW() - $2::INTERVAL
		ORDER BY date ASC
	`

	rows, err := r.db.Query(ctx, query, stockID, interval)
	if err != nil {
		return nil, fmt.Errorf("stock repository: GetHistory: %w", err)
	}
	defer rows.Close()

	var history []domain.StockHistory
	for rows.Next() {
		var h domain.StockHistory
		if err := rows.Scan(&h.StockID, &h.Date, &h.Open, &h.High, &h.Low, &h.Close, &h.Volume); err != nil {
			return nil, fmt.Errorf("stock repository: GetHistory scan: %w", err)
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

func rangeToInterval(r string) string {
	switch r {
	case "1m":
		return "1 month"
	case "3m":
		return "3 months"
	case "6m":
		return "6 months"
	case "5y":
		return "5 years"
	default:
		return "1 year"
	}
}

// GetByTickers returns multiple stocks by their ticker symbols.
func (r *StockRepository) GetByTickers(ctx context.Context, tickers []string) ([]domain.StockDetail, error) {
	if len(tickers) == 0 {
		return []domain.StockDetail{}, nil
	}

	query := `
		SELECT
			s.id, s.ticker, s.company_name, s.sector, s.industry,
			s.market_cap, s.pe_ratio, s.eps, s.revenue_growth,
			s.ranking_score, s.created_at, s.updated_at,
			COALESCE(p.price, 0),
			COALESCE(p.change_percent, 0),
			COALESCE(p.updated_at, NOW())
		FROM stocks s
		LEFT JOIN stock_prices p ON p.stock_id = s.id
		WHERE s.ticker = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, tickers)
	if err != nil {
		return nil, fmt.Errorf("stock repository: GetByTickers: %w", err)
	}
	defer rows.Close()

	var stocks []domain.StockDetail
	for rows.Next() {
		var sd domain.StockDetail
		if err := rows.Scan(
			&sd.ID, &sd.Ticker, &sd.CompanyName, &sd.Sector, &sd.Industry,
			&sd.MarketCap, &sd.PERatio, &sd.EPS, &sd.RevenueGrowth,
			&sd.RankingScore, &sd.CreatedAt, &sd.UpdatedAt,
			&sd.Price, &sd.ChangePercent, &sd.PriceUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("stock repository: GetByTickers scan: %w", err)
		}
		stocks = append(stocks, sd)
	}

	return stocks, rows.Err()
}

// Search returns stocks matching a query (ticker or company name).
func (r *StockRepository) Search(ctx context.Context, query string, limit int) ([]domain.StockDetail, error) {
	if limit <= 0 {
		limit = 20
	}

	sqlQuery := `
		SELECT
			s.id, s.ticker, s.company_name, s.sector, s.industry,
			s.market_cap, s.pe_ratio, s.eps, s.revenue_growth,
			s.ranking_score, s.created_at, s.updated_at,
			COALESCE(p.price, 0),
			COALESCE(p.change_percent, 0),
			COALESCE(p.updated_at, NOW())
		FROM stocks s
		LEFT JOIN stock_prices p ON p.stock_id = s.id
		WHERE s.ticker ILIKE $1 OR s.company_name ILIKE $1
		ORDER BY s.market_cap DESC NULLS LAST
		LIMIT $2
	`

	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(ctx, sqlQuery, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("stock repository: Search: %w", err)
	}
	defer rows.Close()

	var stocks []domain.StockDetail
	for rows.Next() {
		var sd domain.StockDetail
		if err := rows.Scan(
			&sd.ID, &sd.Ticker, &sd.CompanyName, &sd.Sector, &sd.Industry,
			&sd.MarketCap, &sd.PERatio, &sd.EPS, &sd.RevenueGrowth,
			&sd.RankingScore, &sd.CreatedAt, &sd.UpdatedAt,
			&sd.Price, &sd.ChangePercent, &sd.PriceUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("stock repository: Search scan: %w", err)
		}
		stocks = append(stocks, sd)
	}

	return stocks, rows.Err()
}
