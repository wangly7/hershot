package simulator

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/wangly7/hershot/services/ingestion-service/internal/event"
)

type EventPublisher interface {
	PublishGameEvent(
		ctx context.Context,
		gameEvent event.GameEvent,
	) error
}

type Simulator struct {
	publisher EventPublisher
	interval  time.Duration

	gameID     string
	homeTeamID string
	awayTeamID string

	sequence int64

	quarter          int
	remainingSeconds int
	homeScore        int
	awayScore        int

	random *rand.Rand
}

func New(
	publisher EventPublisher,
	interval time.Duration,
	gameID string,
	homeTeamID string,
	awayTeamID string,
) *Simulator {
	return &Simulator{
		publisher: publisher,
		interval:  interval,

		gameID:     gameID,
		homeTeamID: homeTeamID,
		awayTeamID: awayTeamID,

		quarter:          1,
		remainingSeconds: 10 * 60,

		random: rand.New(
			rand.NewSource(time.Now().UnixNano()),
		),
	}
}

func (s *Simulator) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("simulation started for game: %s", s.gameID)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			gameEvent, err := s.nextEvent(ctx)
			if err != nil {
				return fmt.Errorf("generate next event: %w", err)
			}

			if err := s.publisher.PublishGameEvent(ctx, gameEvent); err != nil {
				return fmt.Errorf("publish simulated event: %w", &err)
			}

			log.Printf(
				"published event sequence=%d, type=%s, score=%d-%d, clock=%s",
				gameEvent.Sequence,
				gameEvent.EventType,
				gameEvent.HomeScore,
				gameEvent.AwayScore,
				gameEvent.GameClock,
			)

		}
	}
}

func (s *Simulator) nextEvent(ctx context.Context) (event.GameEvent, error) {
	s.sequence++

	elapsed := s.random.Intn(20) + 5
	s.remainingSeconds -= elapsed

	if s.remainingSeconds < 0 {
		s.remainingSeconds = 0
	}

	eventType := s.RandomEventType()
	teamId := s.randomTeamID()

	points := 0
	metadata := make(map[string]any)

	switch eventType {
	case event.EventTypeShotMade:
		points = s.randomShotPoints()
		if teamId == s.homeTeamID {
			s.homeScore += points
		} else {
			s.awayScore += points
		}

		metadata["shotType"] = shotType(points)
	case event.EventTypeShotMissed:
		metadata["shotType"] = shotType(s.randomShotPoints())
	case event.EventTypeRebound:
		if s.random.Intn(2) == 0 {
			metadata["reboundType"] = "OFFENSIVE"
		} else {
			metadata["reboundType"] = "DEFENSIVE"
		}
	case event.EventTypeTurnOver:
		metadata["turnoverType"] = "BAD_PASS"
	case event.EventTypeFoul:
		metadata["foulType"] = "PERSONAL"
	}

	now := time.Now().UTC()

	gameEvent := event.GameEvent{
		EventID:   uuid.NewString(),
		GameID:    s.gameID,
		Sequence:  s.sequence,
		EventType: eventType,

		TeamID:     teamId,
		HomeTeamID: s.homeTeamID,
		AwayTeamID: s.awayTeamID,
		PlayerID:   s.randomPlayerID(teamId),

		Quarter:   s.quarter,
		GameClock: formatClock(s.remainingSeconds),

		Points: points,

		HomeScore: s.homeScore,
		AwayScore: s.awayScore,

		Metadata: metadata,

		OccuredAt:  now,
		ProducedAt: now,
	}

	if s.remainingSeconds == 0 {
		s.advancePeriod()
	}

	return gameEvent, nil
}

func (s *Simulator) RandomEventType() event.EventType {
	value := s.random.Intn(100)

	switch {
	case value < 35:
		return event.EventTypeShotMade
	case value < 55:
		return event.EventTypeShotMissed
	case value < 72:
		return event.EventTypeRebound
	case value < 82:
		return event.EventTypeTurnOver
	case value < 90:
		return event.EventTypeFoul
	case value < 96:
		return event.EventTypeSteal
	default:
		return event.EventTypeTimeout
	}
}

func (s *Simulator) randomTeamID() string {
	if s.random.Intn(2) == 0 {
		return s.homeTeamID
	}

	return s.awayTeamID
}

func (s *Simulator) randomShotPoints() int {
	value := s.random.Intn(100)

	if value < 70 {
		return 2
	}

	return 3
}

func (s *Simulator) randomPlayerID(teamID string) string {
	playerNumber := s.random.Intn(5) + 1

	return fmt.Sprintf("%s-player-%d", teamID, playerNumber)
}

func (s *Simulator) advancePeriod() {
	if s.quarter > 4 {
		return
	}

	s.quarter++
	s.remainingSeconds = 60 * 10
}

func formatClock(totalSeconds int) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func shotType(points int) string {
	switch points {
	case 3:
		return "THREE_POINTS"
	default:
		return "TWO_POINTS"
	}
}
