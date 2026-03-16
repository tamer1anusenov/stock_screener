package handler

import (
	"errors"
	"net/http"
	"stock_screener/internal/service"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type StockHandler struct {
	stockService *service.StockService
}

func NewStockHandler(stockService *service.StockService) *StockHandler {
	return &StockHandler{stockService: stockService}
}

// GET /stocks/discover
func (h *StockHandler) Discover(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromHeader(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing or invalid X-User-ID header")
		return
	}

	stock, err := h.stockService.Discover(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrNoStocksAvailable) {
			respondError(w, http.StatusNoContent, "no more stocks to discover")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to discover stock")
		return
	}

	respond(w, http.StatusOK, stock)
}

// GET /stocks/{ticker}
func (h *StockHandler) GetByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		respondError(w, http.StatusBadRequest, "ticker is required")
		return
	}

	stock, err := h.stockService.GetByTicker(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, service.ErrStockNotFound) {
			respondError(w, http.StatusNotFound, "stock not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get stock")
		return
	}

	respond(w, http.StatusOK, stock)
}

// GET /stocks/{ticker}/history?range=1y
func (h *StockHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		respondError(w, http.StatusBadRequest, "ticker is required")
		return
	}

	rangeStr := r.URL.Query().Get("range")
	if rangeStr == "" {
		rangeStr = "1y"
	}

	history, err := h.stockService.GetHistory(r.Context(), ticker, rangeStr)
	if err != nil {
		if errors.Is(err, service.ErrStockNotFound) {
			respondError(w, http.StatusNotFound, "stock not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get history")
		return
	}

	respond(w, http.StatusOK, history)
}

// GET /stocks/batch?tickers=AAPL,MSFT,GOOG
func (h *StockHandler) GetBatch(w http.ResponseWriter, r *http.Request) {
	tickersParam := r.URL.Query().Get("tickers")
	if tickersParam == "" {
		respondError(w, http.StatusBadRequest, "tickers query param required")
		return
	}

	tickers := strings.Split(tickersParam, ",")
	for i := range tickers {
		tickers[i] = strings.TrimSpace(tickers[i])
	}

	stocks, err := h.stockService.GetByTickers(r.Context(), tickers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get stocks")
		return
	}

	respond(w, http.StatusOK, stocks)
}

// GET /stocks/search?q=apple
func (h *StockHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "q query param required")
		return
	}

	stocks, err := h.stockService.Search(r.Context(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to search stocks")
		return
	}

	respond(w, http.StatusOK, stocks)
}

// userIDFromHeader extracts the user UUID from the X-User-ID request header.
func userIDFromHeader(r *http.Request) (uuid.UUID, error) {
	raw := r.Header.Get("X-User-ID")
	return uuid.Parse(raw)
}
