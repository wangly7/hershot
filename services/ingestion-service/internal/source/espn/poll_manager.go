package espn

import (
	"context"
	"sync"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

const (
	defaultPollInterval  = 2 * time.Second
	defaultStartLeadTime = 30 * time.Second
)

type runningPoller struct {
	cancel context.CancelFunc
}

type scheduledGame struct {
	game  GameInfo
	timer *time.Timer
}

type PollManager struct {
	client Client

	pollInterval  time.Duration
	startLeadTime time.Duration
	output        chan<- domain.RawPlay

	mu sync.Mutex

	schedules map[string]*scheduledGame
	pollers   map[string]*runningPoller
}

func NewPollManager(
	client Client,
	pollInterval time.Duration,
	startLeadTime time.Duration,
	output chan<- domain.RawPlay,
) *PollManager {
	if pollInterval <= 0 {
		pollInterval = defaultPollInterval
	}
	if startLeadTime <= 0 {
		startLeadTime = defaultStartLeadTime
	}

	return &PollManager{
		client:        client,
		pollInterval:  pollInterval,
		startLeadTime: startLeadTime,
		output:        output,
		schedules:     make(map[string]*scheduledGame),
		pollers:       make(map[string]*runningPoller),
	}
}

func (m *PollManager) Start(
	ctx context.Context,
	game GameInfo,
) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.startLocked(ctx, game)
}

// Schdules creates or updates the timer for one game.
//
// It returns false when the game can not be scheduled.
func (m *PollManager) Schedule(
	ctx context.Context,
	game GameInfo,
) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if game.EventID == "" ||
		game.Completed ||
		game.StartTime.IsZero() {
		return false
	}

	if _, running := m.pollers[game.EventID]; running {
		return false
	}

	if existing, scheduled := m.schedules[game.EventID]; scheduled {
		// if game time does not change
		if sameSchedule(existing.game, game) {
			return false
		}

		// TODO: if game time changed, delete previous timer

	}

	return true
}

func sameSchedule(left GameInfo, right GameInfo) bool {
	return left.EventID == right.EventID &&
		left.CompetitionID == right.CompetitionID &&
		left.StartTime.Equal(right.StartTime)
}

func (m *PollManager) scheduledLocked(
	parentCtx context.Context,
	game GameInfo,
) {
	startAt := game.StartTime.Add(-m.startLeadTime)
	delay := time.Until(startAt)

	if delay < 0 {
		delay = 0
	}

	entry := &scheduledGame{
		game: game,
	}

	entry.timer = time.AfterFunc(delay, func() {
		m.startScheduledPoller(parentCtx, entry)

	})

	m.schedules[game.EventID] = entry
}

func (m *PollManager) startScheduledPoller(
	parentCtx context.Context,
	expected *scheduledGame,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	eventID := expected.game.EventID

	current, exists := m.schedules[eventID]
	if !exists {
		return
	}

	if current != expected {
		return
	}

	delete(m.schedules, eventID)
}

// Start a poller without acquiring m.mu.
//
// The caller must hold m.mu
func (m *PollManager) startLocked(
	parentCtx context.Context,
	game GameInfo,
) bool {
	if game.EventID == "" {
		return false
	}

	if parentCtx.Err() != nil {
		return false
	}

	if _, exists := m.pollers[game.EventID]; exists {
		return false
	}

	pollerCtx, cancel := context.WithCancel(parentCtx)

	running := &runningPoller{
		cancel: cancel,
	}

	poller := NewPoller(
		m.client,
		game,
		m.pollInterval,
		m.output,
	)

	m.pollers[game.EventID] = running

	go func() {
		poller.Run(pollerCtx)
	}()

	return true
}

func (m *PollManager) removeRunningPoller(
	eventID string,
	expected *runningPoller,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, exists := m.pollers[eventID]
	if !exists {
		return
	}

	if current != expected {
		return
	}

	delete(m.pollers, eventID)
}

func (m *PollManager) Stop(eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.stopLocked(eventID)
}

func (m *PollManager) stopLocked(
	eventID string,
) bool {
	if eventID == "" {
		return false
	}
	cancel, exists := m.pollers[eventID]
	if !exists {
		return false
	}

	cancel()
	delete(m.pollers, eventID)

	return true
}

// Sync synchronizes the running Pollers with the supplied game list.
//
// Example:
// Running Pollers: A, B
// Supplied Games: B, C
//
// Sync will:
//   - stop A
//   - keep B
//   - start C
func (m *PollManager) Sync(
	ctx context.Context,
	games []GameInfo,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	desiredGames := make(map[string]GameInfo, len(games))

	for _, game := range games {
		if game.EventID == "" {
			continue
		}
		switch game.State {
		case GameStateIn:
			desiredGames[game.EventID] = game
		default:
			continue
		}
	}

	// stop pollers whose games are no longer present.
	for eventID := range m.pollers {
		if _, exists := desiredGames[eventID]; exists {
			continue
		}

		m.stopLocked(eventID)
	}

	// start pollers for newly discovered games
	for eventID, game := range desiredGames {
		if _, exists := m.pollers[eventID]; exists {
			continue
		}

		m.startLocked(ctx, game)
	}
}

func (m *PollManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for eventID := range m.pollers {
		m.stopLocked(eventID)
	}
}

func (m *PollManager) ActiveCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.pollers)
}

func (m *PollManager) IsRunning(eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.pollers[eventID]
	return exists
}
