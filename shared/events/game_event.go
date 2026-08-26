package events

import "time"

const CurrentSchemaVersion = 1

type EventType string

const (
	EventTypeUnknown EventType = "UNKNOWN"

	EventTypeGameStart   EventType = "GAME_START"
	EventTypeGameEnd     EventType = "GAME_END"
	EventTypePeriodStart EventType = "PERIOD_START"
	EventTypePeriodEnd   EventType = "PERIOD_END"

	EventTypeShotMade   EventType = "SHOT_MADE"
	EventTypeShotMissed EventType = "SHOT_MISSED"

	EventTypeFreeThrowMade   EventType = "FREE_THROW_MADE"
	EventTypeFreeThrowMissed EventType = "FREE_THROW_MISSED"

	EventTypeRebound      EventType = "REBOUND"
	EventTypeTurnover     EventType = "TURNOVER"
	EventTypeSteal        EventType = "STEAL"
	EventTypeFoul         EventType = "FOUL"
	EventTypeTimeout      EventType = "TIMEOUT"
	EventTypeSubstitution EventType = "SUBSTITUTION"
	EventTypeJumpBall     EventType = "JUMP_BALL"
)

type GameEvent struct {
	SchemaVersion int `json:"schemaVersion"`

	EventID string `json:"eventId"`
	GameID  string `json:"gameId"`

	StreamID   string     `json:"streamId"`
	StreamMode StreamMode `json:"streamMode"`

	Source   string `json:"source"`
	Sequence int64  `json:"sequence"`

	EventTypeID      string    `json:"eventTypeId,omitempty"`
	EventType        EventType `json:"eventType"`
	EventTypeName    string    `json:"eventTypeName,omitempty"`
	Description      string    `json:"description,omitempty"`
	ShortDescription string    `json:"shortDescription,omitempty"`

	Period           int    `json:"period"`
	Clock            string `json:"clock"`
	RemainingSeconds int    `json:"remainingSeconds"`

	HomeScore  int `json:"homeScore"`
	AwayScore  int `json:"awayScore"`
	ScoreValue int `jons:"scoreValue"`

	TeamID     string `json:"teamId,omitempty"`
	HomeTeamID string `json:"homeTeamId,omitempty"`
	AwayTeamID string `json:"awayTeamId,omitempty"`

	PlayerID        string `json:"playerId,omitempty"`
	RelatedPlayerID string `json:"relatedPlayerId,omitempty"`

	OccurredAt time.Time `json:"occurredAt"`
	ProducedAt time.Time `json:"producedAt"`
}
