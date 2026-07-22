package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/config"
	"github.com/wangly7/hershot/services/ingestion-service/internal/producer"
	"github.com/wangly7/hershot/services/ingestion-service/internal/simulator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	gameProducer, err := producer.New(
		cfg.RedpandaBrokers,
		cfg.GameEventsTopic,
	)
	if err != nil {
		log.Fatalf("create producer: %v", err)
	}
	defer gameProducer.Close(ctx)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := gameProducer.Ping(pingCtx); err != nil {
		log.Fatalf("connect to Redpanda: %v", err)
	}

	log.Printf(
		"connected to Redpanda brokers=%v topic=%s",
		cfg.RedpandaBrokers,
		cfg.GameEventsTopic,
	)

	gameSimulator := simulator.New(
		gameProducer,
		time.Duration(cfg.SimulationIntervalSeconds)*time.Second,
		cfg.GameID,
		cfg.HomeTeamID,
		cfg.AwayTeamID,
	)

	if err := gameSimulator.Run(ctx); err != nil {
		log.Fatalf("run simulator: %v", err)
	}

	log.Println("ingestion-service stopped")

}
