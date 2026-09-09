package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wangly7/hershot/services/live-game-service/config"
	"github.com/wangly7/hershot/services/live-game-service/internal/repository"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	eventRepository, err := repository.NewDynamoDBRepository(
		ctx,
		cfg.DynamoDBEndpoint,
		cfg.AWSRegion,
		cfg.DynamoDBGameEventsTable,
	)
	if err != nil {
		log.Fatalf(
			"create DynamoDB repository: %v",
			err,
		)
	}

	if err := eventRepository.EnsureTable(ctx); err != nil {
		log.Fatalf(
			"ensure DynamoDB table: %v",
			err,
		)
	}

	log.Printf(
		"DynamoDB table %q is ready",
		cfg.DynamoDBGameEventsTable,
	)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("live-game-service ok"))
	})

	port := fmt.Sprintf(":%d", cfg.HTTPPort)

	log.Println("live-game-service running on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
