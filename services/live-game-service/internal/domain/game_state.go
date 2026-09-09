package domain

import "time"

type GameState struct {
	GameID string `json:"gameId"`

	HomeTeamID string `json:"homeTeamId"`
	AwayTeamID string `json:"awayTeamId"`

	Quarter int `json:"quarter"`

	RemainingSeconds int `json:"remainingSeconds"`

	HomeScore int `json:"homeScore"`
	AwayScore int `json:"awayScore"`

	PossessionTeamID string `json:"possessionTeamID,omitempty"`

	LastSequence int64     `json:"lastSequence"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
