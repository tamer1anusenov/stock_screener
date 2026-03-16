package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json: "id"`
	Email     string    `json: "email"`
	CreatedAt time.Time `json: "created_at"`
}

type SwipeDirection string

const (
	SwipeLeft  SwipeDirection = "left"
	SwipeRight SwipeDirection = "right"
)

type Swipe struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	StockID   uuid.UUID      `json:"stock_id"`
	Direction SwipeDirection `json:"direction"`
	CreatedAt time.Time      `json:"created_at"`
}

type WatchlistEntry struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	StockID   uuid.UUID      `json:"stock_id"`
	Direction SwipeDirection `json:"direction"`
	CreatedAt time.Time      `json:"created_at"`
}

type WatchlistItem struct {
	WatchlistEntry
	Stock StockDetail `json:"stock"`
}
