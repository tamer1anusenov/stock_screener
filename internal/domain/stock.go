package domain

import (
	"time"

	"github.com/google/uuid"
)

// Stock represents the core stock entity with semi-static fundamentals.
// This is a pure domain model — no DB tags, no HTTP tags.
type Stock struct {
	ID            uuid.UUID `json:"id"`
	Ticker        string    `json:"ticker"`
	CompanyName   string    `json:"company_name"`
	Sector        string    `json:"sector"`
	Industry      string    `json:"industry"`
	MarketCap     int64     `json:"market_cap"`
	PERatio       float64   `json:"pe_ratio"`
	EPS           float64   `json:"eps"`
	RevenueGrowth float64   `json:"revenue_growth"`
	RankingScore  float64   `json:"ranking_score"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// StockPrice holds the latest price snapshot for a stock.
type StockPrice struct {
	StockID       uuid.UUID `json:"stock_id"`
	Price         float64   `json:"price"`
	ChangePercent float64   `json:"change_percent"`
	Volume        int64     `json:"volume"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// StockHistory holds a single OHLCV candle for chart rendering.
type StockHistory struct {
	StockID uuid.UUID `json:"stock_id"`
	Date    time.Time `json:"date"`
	Open    float64   `json:"open"`
	High    float64   `json:"high"`
	Low     float64   `json:"low"`
	Close   float64   `json:"close"`
	Volume  int64     `json:"volume"`
}

// StockDetail is a composed view returned by the API —
// stock fundamentals + latest price in one response.
type StockDetail struct {
	Stock
	Price          float64   `json:"price"`
	ChangePercent  float64   `json:"change_percent"`
	PriceUpdatedAt time.Time `json:"price_updated_at"`
}
