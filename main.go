package main

import (
	"context"
	"log"
	"net/http"
	"stock_screener/internal/config"
	"stock_screener/internal/database"
	"stock_screener/internal/handler"
	"stock_screener/internal/repository"
	"stock_screener/internal/service"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("Database connected")

	stockRepo := repository.NewStockRepository(db)
	swipeRepo := repository.NewSwipeRepository(db)
	watchlistRepo := repository.NewWatchlistRepository(db)
	syncRepo := repository.NewSyncRepository(db)

	stockService := service.NewStockService(stockRepo)
	swipeService := service.NewSwipeService(swipeRepo, watchlistRepo)
	watchlistService := service.NewWatchlistService(watchlistRepo)

	stockHandler := handler.NewStockHandler(stockService)
	swipeHandler := handler.NewSwipeHandler(swipeService)
	watchlistHandler := handler.NewWatchlistHandler(watchlistService)
	syncHandler := handler.NewSyncHandler(syncRepo)

	router := handler.NewRouter(stockHandler, swipeHandler, watchlistHandler, syncHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
