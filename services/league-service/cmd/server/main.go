package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wangly7/hershot/services/league-service/config"
	"github.com/wangly7/hershot/services/league-service/internal/database"
	"github.com/wangly7/hershot/services/league-service/internal/team"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.PostgresURL())
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer db.Close()

	teamRepository := team.NewPostgresRepository(db)
	teamService := team.NewService(teamRepository)
	teamHandler := team.NewHandler(teamService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("league-service ok"))
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(pingCtx); err != nil {
			http.Error(w, "postgres unavailable", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	mux.HandleFunc("GET /teams", teamHandler.ListTeams)
	mux.HandleFunc("GET /teams/{id}", teamHandler.GetTeam)
	mux.HandleFunc("POST /teams", teamHandler.CreateTeam)
	mux.HandleFunc("PUT /teams/{id}", teamHandler.UpdateTeam)
	mux.HandleFunc("DELETE /teams/{id}", teamHandler.DeleteTeam)

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("league-service running in %s", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("league-service failed: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown fail: %v", err)
	} else {
		log.Printf("league-service stopped gracefully")
	}
}
