package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncRepository struct {
	db *pgxpool.Pool
}

func NewSyncRepository(db *pgxpool.Pool) *SyncRepository {
	return &SyncRepository{db: db}
}

type SyncStatus struct {
	SyncType      string `json:"sync_type"`
	LastSync      any    `json:"last_sync"`
	Status        string `json:"status"`
	RecordsSynced int    `json:"records_synced"`
	ErrorMessage  any    `json:"error_message"`
}

func (r *SyncRepository) GetSyncStatus(ctx context.Context) ([]SyncStatus, error) {
	query := `
		SELECT sync_type, last_sync, status, records_synced, error_message
		FROM sync_status
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("sync repository: GetSyncStatus: %w", err)
	}
	defer rows.Close()

	var statuses []SyncStatus
	for rows.Next() {
		var s SyncStatus
		if err := rows.Scan(&s.SyncType, &s.LastSync, &s.Status, &s.RecordsSynced, &s.ErrorMessage); err != nil {
			return nil, fmt.Errorf("sync repository: scan: %w", err)
		}
		statuses = append(statuses, s)
	}

	return statuses, rows.Err()
}
