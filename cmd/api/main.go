package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/always-tired/crud-subscriptions/docs"
	"github.com/always-tired/crud-subscriptions/internal/config"
	"github.com/always-tired/crud-subscriptions/internal/logger"
	"github.com/always-tired/crud-subscriptions/internal/repository/postgres"
	httptransport "github.com/always-tired/crud-subscriptions/internal/transport/http"
	"github.com/always-tired/crud-subscriptions/internal/usecase"
)

// @title Subscription Aggregator API
// @version 1.0.0
// @description REST API for subscription aggregation service
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		log.Error("db connect", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.NewSubscriptionRepository(pool)
	service := usecase.NewService(repo, log)
	h := httptransport.NewHandler(service, log)

	r := h.Router()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	r.Get("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		log.Info("http server started", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShutdown)
	log.Info("http server stopped")
}
