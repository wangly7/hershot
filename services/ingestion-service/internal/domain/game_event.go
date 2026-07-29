package domain

import "time"

type EventType string

const (
	EventTypeUnkown EventType = "UNKNOWN"

	EventTypeGameStart   EventType = "GAME_START"
	EventTypeGameEnd     EventType = "GAME_END"
	EventTypePeriodStart EventType = "PERIOD_START"
	EventTypePeriodEnd   EventType = "PERIOD_END"

	EventTypeShotMade   EventType = "SHOT_MADE"
	EventTypeShotMissed EventType = "SHOT_MISSED"

	EventTypeFreeThrowMade   EventType = "FREE_THROW_MADE"
	EventTypeFreeThrowMissed EventType = "FREE_THROW_MISSED"

	EventTypeRebound  EventType = "REBOUND"
	EventTypeTurnOver EventType = "TURNOVER"
	EventTypeSteal    EventType = "STEAL"
	EventTypeFoul     EventType = "FOUL"
	EventTypeTimeout  EventType = "Timeout"

	EventTypeSubstitution EventType = "SUBSTITUTION"
	EventTypeJumpBall     EventType = "JUMP_BALL"
)

type GameEvent struct {
	Source string

	EventID string `json:"eventId"`
	GameID  string `json:"gameId"`

	Provider string `json:"provider"`
	Sequence int64  `json:"sequence"`

	EventTypeID   string    `json:"evnetTypeId,omitempty"`
	EventType     EventType `json:"eventType"`
	EventTypeName string    `json:"eventTypeName,omitempty"`
	Description   string    `json:"description,omitempty"`

	Period           int `json:"period"`
	RemainingSeconds int `json:"remainingSeconds"`

	HomeScore int `json:"homeScore"`
	AwayScore int `json:"awayScore"`

	TeamID     string `json:"teamId,omitempty"`
	HomeTeamID string `json:"homeTeamId,omitempty"`
	AwayTeamID string `json:"awayTeamId,omitempty"`

	PlayerID        string `json:"playerId,omitempty"`
	RelatedPlayerID string `json:"realtedPlayerId,omitempty"`

	OccurredAt  time.Time `json:"occurredAt"`
	PublishedAt time.Time `json:"producedAt"`
}
