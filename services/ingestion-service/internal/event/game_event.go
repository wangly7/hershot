package event

import "time"

type EventType string

const (
	EventTypeShotMade        EventType = "SHOT_MADE"
	EventTypeShotMissed      EventType = "SHOT_MISSED"
	EventTypeRebound         EventType = "Rebound"
	EventTypeTurnOver        EventType = "TURNOVER"
	EventTypeSteal           EventType = "STEAL"
	EventTypeFoul            EventType = "Foul"
	EventTypeFreeThrowMade   EventType = "FREE_THROW_MADE"
	EventTypeFreeThrowMissed EventType = "FreeThrowMiss"
	EventTypeTimeout         EventType = "Timeout"
	EventTypePeriodStart     EventType = "PERIOD_START"
	EventTypePeriodEnd       EventType = "PERIOD_END"
	EventTypeGameEnd         EventType = "GAME_END"
)

type GameEvent struct {
	EventID   string    `json:"eventID"`
	GameID    string    `json:"gameID"`
	Sequence  int64     `json:"sequence"`
	EventType EventType `json:"eventType"`

	TeamID     string `json:"teamID"`
	HomeTeamID string `json:"homeTeamID"`
	AwayTeamID string `json:"awayTeamID"`
	PlayerID   string `json:"playerID"`

	RelatedPlayerID string `json:"realtedPlayerID"`

	Quarter   int    `json:"quarter"`
	GameClock string `json:"gameClock"`

	Points int `json:"points,omitempty"`

	HomeScore int `json:"homeScore"`
	AwayScore int `json:"awayScore"`

	Metadata map[string]any `json:"metadata,omitempty"`

	OccuredAt  time.Time `json:"occuredAt"`
	ProducedAt time.Time `json:"produceddAt"`
}
