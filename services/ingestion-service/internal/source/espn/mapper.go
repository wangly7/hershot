package espn

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

// -----------------------------------------------------------------------------
// ScoreboardDto -> Scoreboard
// -----------------------------------------------------------------------------
func MapScoreboard(dto ScoreboardResponse) ([]GameInfo, error) {
	games := make([]GameInfo, 0, len(dto.Events))

	for _, event := range dto.Events {
		game, err := MapScoreboardEvent(event)
		if err != nil {
			return nil, fmt.Errorf(
				"map scoreboard event %q: %w",
				event.ID,
				err,
			)
		}
		games = append(games, game)
	}
	return games, nil
}

// -----------------------------------------------------------------------------
// ScoreboardEventDto -> GameInfo
// -----------------------------------------------------------------------------
func MapScoreboardEvent(dto ScoreboardEventDTO) (GameInfo, error) {
	competition, err := getPrimaryCompetition(dto)
	if err != nil {
		return GameInfo{}, err
	}

	homeTeamID, awayTeamID, err := mapCompetitorsTeamIDs(competition.Competitors)
	if err != nil {
		return GameInfo{}, fmt.Errorf(
			"map competitors for event %s: %w",
			dto.ID,
			err,
		)
	}

	state, err := mapGameState(dto.Status.Type.State)
	if err != nil {
		return GameInfo{}, err
	}

	startTime, err := parseOptionalTime(dto.Date)
	if err != nil {
		return GameInfo{}, fmt.Errorf("parse start time for event %s: %w", dto.ID, err)
	}

	return GameInfo{
		EventID:       dto.ID,
		CompetitionID: competition.ID,

		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,

		State: state,

		StatusName: dto.Status.Type.Name,
		Completed:  dto.Status.Type.Completed,

		StartTime: startTime,
	}, nil
}

// -----------------------------------------------------------------------------
// PlayDto -> Play
// -----------------------------------------------------------------------------
func MapPlay(gameID string, dto PlayDTO) (domain.RawPlay, error) {
	teamID, err := mapTeamID(dto)
	if err != nil {
		return domain.RawPlay{}, fmt.Errorf("map team id for play %s: %w", dto.ID, err)
	}

	playerID, err := mapPrimaryPlayerID(dto.Participants)
	if err != nil {
		return domain.RawPlay{}, fmt.Errorf(
			"map primary player ID for play %s: %w", dto.ID, err,
		)
	}

	relatedPlayerID, err := mapRelatedPlayerID(dto.Participants)
	if err != nil {
		return domain.RawPlay{}, fmt.Errorf(
			"map related player ID for play %s: %w",
			dto.ID,
			err,
		)
	}

	occurredAt, err := mapOccurredTime(dto)
	if err != nil {
		return domain.RawPlay{}, fmt.Errorf(
			"map occurred time for play %s: %w",
			dto.ID,
			err,
		)
	}

	return domain.RawPlay{
		Provider: "espn",

		GameID:   gameID,
		EventID:  dto.ID,
		Sequence: mapSequence(dto),

		EventTypeID:      dto.Type.ID,
		EventTypeName:    mapRawEventType(dto),
		Description:      dto.Text,
		ShortDescription: dto.ShortText,

		Period:           dto.Period.Number,
		Clock:            dto.Clock.DisplayValue,
		RemainingSeconds: int(dto.Clock.Value),

		HomeScore: dto.HomeScore,
		AwayScore: dto.AwayScore,

		TeamID: teamID,

		ScoreValue: dto.ScoreValue,

		PlayerID:        playerID,
		RelatedPlayerID: relatedPlayerID,

		OccurredAt: occurredAt,
	}, nil
}

// -----------------------------------------------------------------------------
// Scoreboard helpers
// -----------------------------------------------------------------------------
func getPrimaryCompetition(dto ScoreboardEventDTO) (CompetitionDTO, error) {
	if len(dto.Competitions) == 0 {
		return CompetitionDTO{}, fmt.Errorf("scoreboard event %s has no competitions", dto.ID)
	}
	competition := dto.Competitions[0]

	return competition, nil
}

func mapCompetitorsTeamIDs(
	competitors []CompetitorDTO,
) (homeTeamID string, awayTeamID string, err error) {
	for _, competitor := range competitors {
		teamID := strings.TrimSpace(competitor.Team.ID)

		switch strings.ToLower(strings.TrimSpace(competitor.HomeAway)) {
		case "home":
			homeTeamID = teamID
		case "away":
			awayTeamID = teamID
		}
	}

	if homeTeamID == "" {
		return "", "", fmt.Errorf("home team ID is missing")
	}

	if awayTeamID == "" {
		return "", "", fmt.Errorf("away team ID is missing")
	}

	return homeTeamID, awayTeamID, nil
}

func mapGameState(state string) (GameState, error) {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "pre":
		return GameStatePre, nil
	case "in":
		return GameStateIn, nil
	case "post":
		return GameStatePost, nil
	default:
		return "", fmt.Errorf("unsupported ESPN game state %q", state)
	}
}

// -----------------------------------------------------------------------------
// Play helpers
// -----------------------------------------------------------------------------
func mapTeamID(dto PlayDTO) (string, error) {
	return extractIDFromRef(dto.Team.Ref)
}

func findParticipantByOrder(participants []ParticipantDTO, order int) (ParticipantDTO, bool) {
	for _, participant := range participants {
		if participant.Order == order {
			return participant, true
		}
	}

	return ParticipantDTO{}, false
}

// ESPN participants usually use order=1 for the main player involved
// in the play
func mapPrimaryPlayerID(participants []ParticipantDTO) (string, error) {
	// TODO: support more scenes
	primary, ok := findParticipantByOrder(participants, 1)
	if !ok {
		return "", nil
	}
	return extractIDFromRef(primary.Athlete.Ref)
}

// The second participant may represent an assister, defender,
// rebound opponent, substituted player, etc
func mapRelatedPlayerID(participants []ParticipantDTO) (string, error) {
	participant, ok := findParticipantByOrder(participants, 2)
	if !ok {
		return "", nil
	}
	return extractIDFromRef(participant.Athlete.Ref)
}

func mapSequence(dto PlayDTO) int64 {
	return int64(dto.SequenceNumber)
}

func mapRawEventType(dto PlayDTO) string {
	if dto.Type.Text != "" {
		return dto.Type.Text
	}

	return dto.Type.ID
}

func mapOccurredTime(dto PlayDTO) (time.Time, error) {
	return parseOptionalTime(dto.Modified)
}

func parseOptionalTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		"2006-01-02T15:04Z",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %q", value)
}

func extractIDFromRef(ref string) (string, error) {
	ref = strings.TrimSpace(ref)

	if ref == "" {
		return "", nil
	}

	parsed, err := url.Parse(ref)
	if err != nil {
		return "", fmt.Errorf("parse ref %q: %w", ref, err)
	}

	id := path.Base(strings.TrimSuffix(parsed.Path, "/"))

	if id == "" || id == "." || id == "/" {
		return "", fmt.Errorf("ref %q does not contain an ID", ref)
	}

	return id, nil
}
