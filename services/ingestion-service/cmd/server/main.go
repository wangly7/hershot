package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/wangly7/hershot/services/ingestion-service/config"
	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
	"github.com/wangly7/hershot/services/ingestion-service/internal/eventmapper"
	"github.com/wangly7/hershot/services/ingestion-service/internal/kafka"
	"github.com/wangly7/hershot/services/ingestion-service/internal/source/espn"
	"github.com/wangly7/hershot/services/ingestion-service/internal/source/simulator"
)

func main() {
	if err := run(); err != nil {
		log.Printf("ingestion service stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// The context is canceld when the process receives Ctrl+C or SIGTERM.
	rootCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	/*
		Create the ESPN API client.
	*/
	var client espn.Client

	switch strings.ToLower(cfg.ClientType) {
	case "espn":
		client = espn.NewClient(
			espn.ClientConfig{
				SiteBaseURL: cfg.ESPNSiteBaseURL,
				CoreBaseURL: cfg.ESPNCoreBaseURL,
				Timeout:     cfg.ESPNHTTPTimeout,
				PlaysLimits: cfg.ESPNPlaysLimits,
			},
		)
	case "simulator":
		client, err = simulator.NewClientFromFixtures("./internal/source/simulator/testdata/")
		if err != nil {
			return fmt.Errorf("create simulator client: %w", err)
		}
	default:
		return fmt.Errorf("unsupported client type: %s", cfg.ClientType)
	}

	eventSource := espn.NewSource(client, espn.SourceConfig{
		PollInterval:            cfg.PollInterval,
		StartLeadTime:           cfg.StartLeadTime,
		ScheduleRefreshInterval: cfg.ScheduleRefreshInterval,
	})

	/*
		Create the Kafka producer
	*/
	eventProducer, err := kafka.NewProducer(
		cfg.KafkaBrokers,
		cfg.KafkaGameEventsTopic,
	)
	if err != nil {
		return fmt.Errorf("create kafka producer: %w", err)
	}
	defer eventProducer.Close()

	output := make(
		chan domain.RawPlay,
		cfg.OutputBuffers,
	)

	/*
		Pipeline: Sheduler -> PollManager -> Pollers
	*/
	runCtx, cancelRun := context.WithCancel(rootCtx)
	defer cancelRun()

	sourceErr := make(chan error, 1)

	go func() {
		sourceErr <- eventSource.Run(runCtx, output)
	}()

	log.Printf(
		"ingestion service started brokers=%v topic=%s",
		cfg.KafkaBrokers,
		cfg.KafkaGameEventsTopic,
	)

	for {
		select {
		case play, ok := <-output:
			if !ok {
				return nil
			}

			event, err := eventmapper.ToGameEvent(play)
			if err != nil {
				return fmt.Errorf(
					"map raw play playEventID=%s gameID=%s: %w",
					play.EventID,
					play.GameID,
					err,
				)
			}

			if err := eventProducer.Publish(runCtx, event); err != nil {
				return fmt.Errorf(
					"publish game event eventID=%s gameID=%s: %w",
					event.EventID,
					event.GameID,
					err,
				)
			}
		case err := <-sourceErr:
			if err == nil {
				return nil
			}

			if errors.Is(err, context.Canceled) ||
				errors.Is(err, context.DeadlineExceeded) {
				return nil
			}

			return fmt.Errorf("run ESPN source: %w", err)
		case <-rootCtx.Done():
			log.Println("shutting down ingestion service")
			return nil
		}
	}
}
