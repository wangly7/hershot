package team

import "errors"

var (
	ErrTeamNotFound            = errors.New("team not found")
	ErrInvalidID               = errors.New("invalid team id")
	ErrInvalidTeamName         = errors.New("invalid team name")
	ErrInvalidTeamCity         = errors.New("invalid team city")
	ErrInvalidTeamAbbreviation = errors.New("invalid team abbreviation")
)
