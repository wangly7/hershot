package eventmapper

import (
	"strings"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
	"github.com/wangly7/hershot/shared/events"
)

func toEventType(play domain.RawPlay) events.EventType {
	typeName := strings.ToLower(
		strings.TrimSpace(play.EventTypeName),
	)

	description := strings.ToLower(
		strings.TrimSpace(play.Description),
	)

	switch {
	case strings.Contains(typeName, "free throw"):
		if strings.Contains(description, "misses") {
			return events.EventTypeFreeThrowMissed
		}

		if play.ScoreValue > 0 ||
			strings.Contains(description, "makes") {
			return events.EventTypeFreeThrowMade
		}

	case strings.Contains(typeName, "rebound"):
		return events.EventTypeRebound

	case strings.Contains(typeName, "turnover"),
		strings.Contains(typeName, "traveling"):
		return events.EventTypeTurnover

	case strings.Contains(typeName, "steal"):
		return events.EventTypeSteal

	case strings.Contains(typeName, "foul"):
		return events.EventTypeFoul

	case strings.Contains(typeName, "timeout"):
		return events.EventTypeTimeout

	case strings.Contains(typeName, "substitution"):
		return events.EventTypeSubstitution

	case strings.Contains(typeName, "jumpball"),
		strings.Contains(typeName, "jump ball"):
		return events.EventTypeJumpBall

	case isShotType(typeName):
		if play.ScoreValue > 0 ||
			strings.Contains(description, "makes") {
			return events.EventTypeShotMade
		}

		if strings.Contains(description, "misses") {
			return events.EventTypeShotMissed
		}
	}

	return events.EventTypeUnknown
}

func isShotType(typeName string) bool {
	shotTokens := []string{
		"shot",
		"jumper",
		"layup",
		"dunk",
		"hook",
	}

	for _, token := range shotTokens {
		if strings.Contains(typeName, token) {
			return true
		}
	}

	return false
}
