package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"stock_screener/internal/domain"
	"stock_screener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SwipeHandler struct {
	swipeService *service.SwipeService
}

func NewSwipeHandler(swipeService *service.SwipeService) *SwipeHandler {
	return &SwipeHandler{swipeService: swipeService}
}

type swipeRequest struct {
	Direction string `json:"direction"`
}

// POST /stocks/{id}/swipe
// Body: { "direction": "left" | "right" }
// Right swipe automatically adds the stock to the watchlist.
func (h *SwipeHandler) Swipe(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromHeader(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing or invalid X-User-ID header")
		return
	}

	stockID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stock id")
		return
	}

	var req swipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	direction := domain.SwipeDirection(req.Direction)
	if direction != domain.SwipeLeft && direction != domain.SwipeRight {
		respondError(w, http.StatusBadRequest, `direction must be "left" or "right"`)
		return
	}

	swipe, err := h.swipeService.Swipe(r.Context(), userID, stockID, direction)
	if err != nil {
		if errors.Is(err, service.ErrAlreadySwiped) {
			respondError(w, http.StatusConflict, "already swiped this stock")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to record swipe")
		return
	}

	respond(w, http.StatusCreated, swipe)
}
