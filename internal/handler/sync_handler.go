package handler

import (
	"net/http"
	"stock_screener/internal/repository"
)

type SyncHandler struct {
	syncRepo *repository.SyncRepository
}

func NewSyncHandler(syncRepo *repository.SyncRepository) *SyncHandler {
	return &SyncHandler{syncRepo: syncRepo}
}

func (h *SyncHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.syncRepo.GetSyncStatus(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get sync status")
		return
	}

	respond(w, http.StatusOK, map[string]interface{}{
		"syncs": statuses,
	})
}
