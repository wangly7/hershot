package simulator

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/source/espn"
)

func LoadScoreboardFixture(
	path string,
) (espn.ScoreboardResponse, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return espn.ScoreboardResponse{}, fmt.Errorf(
			"read scoreboard fixture %q: %w",
			path,
			err,
		)
	}

	var response espn.ScoreboardResponse

	if err := json.Unmarshal(data, &response); err != nil {
		return espn.ScoreboardResponse{}, fmt.Errorf(
			"decode scoreboard fixture %q: %w",
			path,
			err,
		)
	}

	return response, nil
}

func LoadPlaysFixture(
	path string,
) (espn.PlayResponse, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return espn.PlayResponse{}, fmt.Errorf(
			"read plays fixture %q: %w",
			path,
			err,
		)
	}

	var response espn.PlayResponse

	if err := json.Unmarshal(data, &response); err != nil {
		return espn.PlayResponse{}, fmt.Errorf(
			"decode plays fixture %q: %w",
			path,
			err,
		)
	}

	return response, nil
}

func NewClientFromFixtures(
	path string,
) (*client, error) {
	scoreboard, err := LoadScoreboardFixture(path + "scoreboard.json")
	if err != nil {
		return nil, err
	}
	for i := range scoreboard.Events {
		scoreboard.Events[i].Date = time.Now().
			Add(10 * time.Second).
			UTC().
			Format(time.RFC3339)

		scoreboard.Events[i].Status.Type.State = "in"
		scoreboard.Events[i].Status.Type.Completed = false
	}
	client := NewClient(scoreboard)

	for _, game := range scoreboard.Events {
		eventID := game.ID
		competitionID := game.Competitions[0].ID

		plays, err := LoadPlaysFixture(path + "/" + fmt.Sprintf("%s-%s.json", eventID, competitionID))
		if err != nil {
			return nil, err
		}

		client.SetPlays(
			eventID,
			competitionID,
			plays,
		)
	}
	return client, nil
}
