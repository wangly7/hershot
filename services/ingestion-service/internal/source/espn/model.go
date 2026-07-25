package espn

import "time"

type GameState string

const (
	GameStatePre  GameState = "pre"
	GameStateIn   GameState = "in"
	GameStatePost GameState = "post"
)

type GameInfo struct {
	EventID       string
	CompetitionID string

	HomeTeamID string
	AwayTeamID string

	State GameState

	StatusName string
	Completed  bool

	StartTime time.Time
}
