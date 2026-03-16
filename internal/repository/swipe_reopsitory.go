package repository

import (
	"context"
	"fmt"
	"stock_screener/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SwipeRepository struct {
	db *pgxpool.Pool
}

func NewSwipeRepository(db *pgxpool.Pool) *SwipeRepository {
	return &SwipeRepository{db: db}
}

// Create inserts a swipe record.
// Returns an error if the user has already swiped this stock (unique constraint).
func (r *SwipeRepository) Create(ctx context.Context, userID, stockID uuid.UUID, direction domain.SwipeDirection) (*domain.Swipe, error) {
	query := `
		INSERT INTO swipes (user_id, stock_id, direction)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, stock_id, direction, created_at
	`

	var swipe domain.Swipe
	err := r.db.QueryRow(ctx, query, userID, stockID, direction).Scan(
		&swipe.ID, &swipe.UserID, &swipe.StockID, &swipe.Direction, &swipe.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("swipe repository: Create: %w", err)
	}

	return &swipe, nil
}

// HasSwiped returns true if the user has already swiped this stock.
func (r *SwipeRepository) HasSwiped(ctx context.Context, userID, stockID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM swipes WHERE user_id = $1 AND stock_id = $2)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, userID, stockID).Scan(&exists); err != nil {
		return false, fmt.Errorf("swipe repository: HasSwiped: %w", err)
	}

	return exists, nil
}
