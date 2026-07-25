package espn

import (
	"encoding/json"
	"os"
	"testing"
)

func TestUnmarshalPlay(t *testing.T) {
	data, err := os.ReadFile("../../../testdata/play.json")
	if err != nil {
		t.Fatalf("failed to read play.json")
	}

	var play PlayDTO

	if err := json.Unmarshal(data, &play); err != nil {
		t.Fatalf("failed to unmarshall play: %v", err)
	}

	t.Log("========== Play ==========")
	t.Logf("ID: %s", play.ID)
	t.Logf("Sequence: %d", play.SequenceNumber)
	t.Logf("Type: %s (%s)", play.Type.ID, play.Type.Text)
	t.Logf("Text: %s", play.Text)
	t.Logf("Period: %d", play.Period.Number)
	t.Logf("Clock: %.0f", play.Clock.Value)
	t.Logf("Home Score: %d", play.HomeScore)
	t.Logf("Away Score: %d", play.AwayScore)
	t.Logf("Team Ref: %s", play.Team.Ref)
	t.Logf("Participants: %d", len(play.Participants))
}

func TestUnmarshalPlaysResponse(t *testing.T) {
	data, err := os.ReadFile("../../../testdata/plays.json")
	if err != nil {
		t.Fatalf("failed to read plays.json")
	}

	var resp PlayResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshall play: %v", err)
	}

	t.Logf("Count: %d", resp.Count)
	t.Logf("PageIndex: %d", resp.PageIndex)
	t.Logf("PageSize: %d", resp.PageSize)
	t.Logf("PageCount: %d", resp.PageCount)
	t.Logf("Items: %d", len(resp.Items))
}

func TestUnmarshalScoreboardResponse(t *testing.T) {
	data, err := os.ReadFile("../../../testdata/scoreboard.json")
	if err != nil {
		t.Fatalf("failed to read scoreboard.json")
	}

	var resp ScoreboardResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshall play: %v", err)
	}

	event := resp.Events[2]

	t.Logf("EventID: %s", event.ID)

	comp := event.Competitions[0]

	t.Logf("CompetitionID: %s", comp.ID)

	for _, c := range comp.Competitors {
		t.Logf(
			"Team=%s Score=%s HomeAway=%s",
			c.Team.ID,
			c.Score,
			c.HomeAway,
		)
	}
}
