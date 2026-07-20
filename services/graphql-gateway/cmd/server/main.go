package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/wangly7/hershot/services/graphql-gateway/config"
	"github.com/wangly7/hershot/services/graphql-gateway/internal/graph/generated"
	"github.com/wangly7/hershot/services/graphql-gateway/internal/leagueclient"
	"github.com/wangly7/hershot/services/graphql-gateway/internal/resolver"
)

func main() {
	cfg := config.Load()

	leagueClient := leagueclient.New(cfg.LeagueServiceURL)

	resolver := &resolver.Resolver{
		LeagueClient: leagueClient,
	}

	graphqlServer := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	mux := http.NewServeMux()

	mux.Handle(
		"/graphql",
		graphqlServer,
	)

	mux.Handle(
		"/",
		playground.Handler(
			"GraphQL Playgroud",
			"/graphql",
		),
	)

	mux.HandleFunc("GET /healthz", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("graphql-gateway ok"))
	})

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("graphql-gateway running on %s", addr)
	log.Fatal(server.ListenAndServe())
}
