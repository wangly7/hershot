package eventmapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

// ToGameEvent converts an internal RawPlay into the event contract
// published by ingestion-service.
func ToGameEvent(
	play domain.RawPlay,
) (domain.GameEvent, error) {
	if play.EventID == "" {
		return domain.GameEvent{}, errors.New("raw play event ID is empty")
	}

	if play.GameID == "" {
		return domain.GameEvent{}, errors.New("raw play game ID is empty")
	}

	if play.Sequence <= 0 {
		return domain.GameEvent{}, fmt.Errorf(
			"raw play sequence must be positive: %d",
			play.Sequence,
		)
	}

	return domain.GameEvent{
		EventID: play.EventID,

		GameID: play.GameID,

		Sequence: play.Sequence,

		EventTypeID:   play.EventTypeID,
		EventTypeName: play.EventTypeName,

		Description: play.Description,

		Period:           play.Period,
		RemainingSeconds: play.RemainingSeconds,

		HomeScore: play.HomeScore,
		AwayScore: play.AwayScore,

		TeamID:          play.TeamID,
		PlayerID:        play.PlayerID,
		RelatedPlayerID: play.RelatedPlayerID,

		OccurredAt:  play.OccurredAt,
		PublishedAt: time.Now(),

		Source: play.Provider,
	}, nil
}
