package simulator

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/source/espn"
)

type gameKey struct {
	eventID       string
	competitionID string
}

type client struct {
	mu sync.RWMutex

	scoreboard espn.ScoreboardResponse
	plays      map[gameKey]espn.PlayResponse

	scoreboardErr error
	playsErr      error
}

var ErrPlaysNotFound = errors.New("simulated plays not found")

func NewClient(
	scoreboard espn.ScoreboardResponse,
) *client {
	return &client{
		scoreboard: scoreboard,
		plays:      make(map[gameKey]espn.PlayResponse),
	}
}

func (c *client) GetScoreboard(
	ctx context.Context,
	date time.Time,
) (espn.ScoreboardResponse, error) {
	if err := ctx.Err(); err != nil {
		return espn.ScoreboardResponse{}, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.scoreboardErr != nil {
		return espn.ScoreboardResponse{}, c.scoreboardErr
	}

	return c.scoreboard, nil

}

func (c *client) GetPlays(
	ctx context.Context,
	eventID string,
	competitionID string,
) (espn.PlayResponse, error) {
	if err := ctx.Err(); err != nil {
		return espn.PlayResponse{}, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.playsErr != nil {
		return espn.PlayResponse{}, c.playsErr
	}

	key := gameKey{
		eventID:       eventID,
		competitionID: competitionID,
	}

	response, exists := c.plays[key]
	if !exists {
		return espn.PlayResponse{}, fmt.Errorf(
			"%w: eventID=%s competitionID=%s",
			ErrPlaysNotFound,
			eventID,
			competitionID,
		)
	}

	return response, nil
}

func (c *client) SetScoreboard(
	scoreboard espn.ScoreboardResponse,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.scoreboard = scoreboard
}

func (c *client) SetPlays(
	eventID string,
	competitionID string,
	plays espn.PlayResponse,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := gameKey{
		eventID:       eventID,
		competitionID: competitionID,
	}

	c.plays[key] = plays
}

func (c *client) SetScoreboardError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.scoreboardErr = err
}

func (c *client) SetPlaysError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.playsErr = err
}
