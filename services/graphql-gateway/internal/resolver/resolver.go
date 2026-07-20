package resolver

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import "github.com/wangly7/hershot/services/graphql-gateway/internal/leagueclient"

type Resolver struct {
	LeagueClient *leagueclient.Client
}
