package domain

import (
	"time"

	"github.com/wangly7/hershot/shared/events"
)

type RawPlay struct {
	Source string

	StreamID   string
	StreamMode events.StreamMode

	GameID   string
	EventID  string
	Sequence int64

	EventTypeID      string
	EventTypeName    string
	Description      string
	ShortDescription string

	Period           int
	Clock            string
	RemainingSeconds int

	HomeScore int
	AwayScore int

	TeamID string

	ScoreValue int

	PlayerID        string
	RelatedPlayerID string

	OccurredAt time.Time
}
