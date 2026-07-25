package simulator

import (
	"context"
	"fmt"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

type Source struct {
	gameID   string
	interval time.Duration
}

func New(
	gameID string,
	interval time.Duration,
) *Source {
	return &Source{
		gameID:   gameID,
		interval: interval,
	}
}

func (s *Source) Run(
	ctx context.Context,
	output chan<- domain.RawPlay,
) error {
	plays := s.demoPlays()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for _, play := range plays {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			select {
			case output <- play:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

func (s *Source) demoPlays() []domain.RawPlay {
	now := time.Now().UTC()

	return []domain.RawPlay{
		{
			Provider:         "simulator",
			GameID:           s.gameID,
			EventID:          fmt.Sprintf("%s-%d", s.gameID, 1),
			Sequence:         1,
			EventType:        "period-start",
			Description:      "Start of first quarter",
			Quarter:          1,
			RemainingSeconds: 600,
			HomeScore:        0,
			AwayScore:        0,
			SourceHomeTeamID: "valkyries",
			SourceAwayTeamID: "mystics",
			OccurredAt:       now,
		},
		{
			Provider:         "simulator",
			GameID:           s.gameID,
			EventID:          fmt.Sprintf("%s-%d", s.gameID, 2),
			Sequence:         2,
			EventType:        "made-shot",
			Description:      "Valkyries make a two-point shot",
			Quarter:          1,
			RemainingSeconds: 582,
			HomeScore:        2,
			AwayScore:        0,
			SourceHomeTeamID: "valkyries",
			SourceAwayTeamID: "mystics",
			SourceTeamID:     "valkyries",
			PlayerID:         "player-001",
			OccurredAt:       now.Add(18 * time.Second),
		},
		{
			Provider:         "simulator",
			GameID:           s.gameID,
			EventID:          fmt.Sprintf("%s-%d", s.gameID, 3),
			Sequence:         3,
			EventType:        "turnover",
			Description:      "Mystics turnover",
			Quarter:          1,
			RemainingSeconds: 560,
			HomeScore:        2,
			AwayScore:        0,
			SourceHomeTeamID: "valkyries",
			SourceAwayTeamID: "mystics",
			SourceTeamID:     "mystics",
			PlayerID:         "player-002",
			OccurredAt:       now.Add(40 * time.Second),
		},
	}
}
