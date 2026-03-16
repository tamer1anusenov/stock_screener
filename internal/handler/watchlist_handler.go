package handler

import (
	"net/http"
	"stock_screener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WatchlistHandler struct {
	watchlistService *service.WatchlistService
}

func NewWatchlistHandler(watchlistService *service.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{watchlistService: watchlistService}
}

// GET /watchlist
// Returns all watchlist items for the current user.
func (h *WatchlistHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromHeader(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing or invalid X-User-ID header")
		return
	}

	items, err := h.watchlistService.GetWatchlist(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get watchlist")
		return
	}

	respond(w, http.StatusOK, items)
}

// DELETE /watchlist/{stockId}
// Removes a stock from the user's watchlist.
func (h *WatchlistHandler) Remove(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromHeader(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing or invalid X-User-ID header")
		return
	}

	stockID, err := uuid.Parse(chi.URLParam(r, "stockId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stock id")
		return
	}

	if err := h.watchlistService.Remove(r.Context(), userID, stockID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to remove from watchlist")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
