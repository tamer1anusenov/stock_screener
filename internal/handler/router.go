package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter wires all handlers to their routes and returns the root http.Handler.
//
// Route map:
//
//	GET  /health                      → health check
//	GET  /sync/status                 → worker sync status
//	GET  /stocks/discover             → next stock to swipe
//	GET  /stocks/batch                → bulk stocks (?tickers=AAPL,MSFT)
//	GET  /stocks/search               → search stocks (?q=apple)
//	GET  /stocks/{ticker}             → stock detail + price
//	GET  /stocks/{ticker}/history     → OHLCV time-series (?range=1y)
//	POST /stocks/{id}/swipe           → record a swipe
//	GET  /watchlist                   → get user's watchlist
//	DELETE /watchlist/{stockId}       → remove from watchlist
func NewRouter(
	stockHandler *StockHandler,
	swipeHandler *SwipeHandler,
	watchlistHandler *WatchlistHandler,
	syncHandler *SyncHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/sync/status", syncHandler.GetStatus)

	r.Get("/stocks/discover", stockHandler.Discover)
	r.Get("/stocks/batch", stockHandler.GetBatch)
	r.Get("/stocks/search", stockHandler.Search)
	r.Get("/stocks/{ticker}", stockHandler.GetByTicker)
	r.Get("/stocks/{ticker}/history", stockHandler.GetHistory)

	r.Post("/stocks/{id}/swipe", swipeHandler.Swipe)

	r.Get("/watchlist", watchlistHandler.GetWatchlist)
	r.Delete("/watchlist/{stockId}", watchlistHandler.Remove)

	return r
}
