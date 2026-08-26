package eventmapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
	"github.com/wangly7/hershot/shared/events"
)

// ToGameEvent converts an internal RawPlay into the event contract
// published by ingestion-service.
func ToGameEvent(play domain.RawPlay) (events.GameEvent, error) {
	if play.EventID == "" {
		return events.GameEvent{}, errors.New("raw play event ID is empty")
	}

	if play.GameID == "" {
		return events.GameEvent{}, errors.New("raw play game ID is empty")
	}

	if play.StreamID == "" {
		return events.GameEvent{}, errors.New("raw play stream ID is empty")
	}

	if play.Source == "" {
		return events.GameEvent{}, errors.New("raw play source is empty")
	}

	if play.Sequence <= 0 {
		return events.GameEvent{}, fmt.Errorf(
			"raw play sequence must be positive: %d",
			play.Sequence,
		)
	}

	if err := validateStreamMode(play.StreamMode); err != nil {
		return events.GameEvent{}, err
	}

	return events.GameEvent{
		SchemaVersion: events.CurrentSchemaVersion,

		EventID: play.EventID,
		GameID:  play.GameID,

		StreamID:   play.StreamID,
		StreamMode: play.StreamMode,
		Source:     play.Source,

		Sequence: play.Sequence,

		EventTypeID:   play.EventTypeID,
		EventType:     toEventType(play),
		EventTypeName: play.EventTypeName,

		Description:      play.Description,
		ShortDescription: play.ShortDescription,

		Period:           play.Period,
		Clock:            play.Clock,
		RemainingSeconds: play.RemainingSeconds,

		HomeScore:  play.HomeScore,
		AwayScore:  play.AwayScore,
		ScoreValue: play.ScoreValue,

		TeamID: play.TeamID,

		PlayerID:        play.PlayerID,
		RelatedPlayerID: play.RelatedPlayerID,

		OccurredAt: play.OccurredAt,
		ProducedAt: time.Now().UTC(),
	}, nil
}

func validateStreamMode(mode events.StreamMode) error {
	switch mode {
	case events.StreamModeLive, events.StreamModeReplay:
		return nil
	default:
		return fmt.Errorf("invalid stream mode: %q", mode)
	}
}
