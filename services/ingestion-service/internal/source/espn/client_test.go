package espn

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetScoreboard(t *testing.T) {
	client := NewClient(ClientConfig{})

	date := time.Date(
		2026,
		time.June,
		27,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	resp, err := client.GetScoreboard(context.Background(), date)
	if err != nil {
		t.Fatal("fail to get scoreboard")
	}

	t.Logf("games in total: %d", len(resp.Events))
}

func TestClient_GetPlays(t *testing.T) {
	client := NewClient(ClientConfig{})

	eventID := "401857026"
	competitionID := "401857026"

	resp, err := client.GetPlays(context.Background(), eventID, competitionID)
	if err != nil {
		t.Fatal("fail to get plays")
	}

	t.Logf("game events in total : %d", len(resp.Items))
}
