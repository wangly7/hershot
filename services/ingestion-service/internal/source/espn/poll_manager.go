package espn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
	"github.com/wangly7/hershot/shared/events"
)

const (
	defaultPollInterval  = 2 * time.Second
	defaultStartLeadTime = 30 * time.Second
)

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
	pollers   map[string]context.CancelFunc
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
	if startLeadTime < 0 {
		startLeadTime = defaultStartLeadTime
	}

	return &PollManager{
		client:        client,
		pollInterval:  pollInterval,
		startLeadTime: startLeadTime,
		output:        output,
		schedules:     make(map[string]*scheduledGame),
		pollers:       make(map[string]context.CancelFunc),
	}
}

// SyncSchedule synchronize timers with the latest scroreboard.
//
// It:
//   - creates timers for newly discovered games
//   - replace timers when StartTime changes
//   - remove timers for games no longer present
//   - stops pollers for games marked completed
func (m *PollManager) SyncSchedule(
	ctx context.Context,
	games []GameInfo,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	allGames := make(map[string]GameInfo, len(games))
	schedulableGames := make(map[string]GameInfo, len(games))
	for _, game := range games {
		if game.EventID == "" {
			continue
		}

		allGames[game.EventID] = game

		if game.Completed {
			continue
		}

		if game.StartTime.IsZero() {
			continue
		}

		schedulableGames[game.EventID] = game
	}

	// Remove timers for games that are no longer present or schedule
	for eventID := range m.schedules {
		if _, exists := schedulableGames[eventID]; !exists {
			m.cancelScheduleLocked(eventID)
		}
	}

	// stop timers and poller for game explicitly marked completed.
	for eventID, game := range allGames {
		if game.Completed {
			m.cancelScheduleLocked(eventID)
			m.stopLocked(eventID)
		}
	}

	// create or update schedules
	for eventID, game := range schedulableGames {
		// already started polling
		if _, exists := m.pollers[eventID]; exists {
			continue
		}
		existing, scheduled := m.schedules[eventID]
		if scheduled {
			if sameSchedule(existing.game, game) {
				continue
			}
			m.cancelScheduleLocked(eventID)
		}
		m.scheduleLocked(ctx, game)
	}
}

func (m *PollManager) Start(
	ctx context.Context,
	game GameInfo,
) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cancelScheduleLocked(game.EventID)

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

		// if game time changed, delete previous timer
		m.cancelScheduleLocked(game.EventID)
	}

	m.scheduleLocked(ctx, game)
	return true
}

func sameSchedule(left GameInfo, right GameInfo) bool {
	return left.EventID == right.EventID &&
		left.CompetitionID == right.CompetitionID &&
		left.StartTime.Equal(right.StartTime)
}

func (m *PollManager) scheduleLocked(
	parentCtx context.Context,
	game GameInfo,
) {
	startAt := game.StartTime.Add(-m.startLeadTime)
	delay := time.Until(startAt)
	fmt.Printf("competition %s game delay: %v\n", game.EventID, delay)

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

// startSchedulePoller is called when  a game timer expires.
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

	// ignore an old timer callback after the schedule was replaced
	if current != expected {
		return
	}

	delete(m.schedules, eventID)
	m.startLocked(parentCtx, expected.game)
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

	streamId := "live:" + game.EventID

	poller := NewPoller(
		m.client,
		game,
		streamId,
		events.StreamModeLive,
		m.pollInterval,
		m.output,
	)

	m.pollers[game.EventID] = cancel

	go poller.Run(pollerCtx)

	return true
}

// Stop stops a running poller and cancels any pending timer.
func (m *PollManager) Stop(eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	scheduleCanceled := m.cancelScheduleLocked(eventID)
	pollerStopped := m.stopLocked(eventID)

	return scheduleCanceled || pollerStopped
}

// stopLocked stops a running poller.
//
// The caller must already hold m.mu.
func (m *PollManager) stopLocked(
	eventID string,
) bool {
	cancel, exists := m.pollers[eventID]
	if !exists {
		return false
	}

	// delete before cancel
	delete(m.pollers, eventID)
	cancel()

	return true
}

// cancelScheduleLocked cancels one pending timer
//
// The caller must already hold m.mu.
func (m *PollManager) cancelScheduleLocked(eventID string) bool {
	scheduled, exists := m.schedules[eventID]
	if !exists {
		return false
	}

	delete(m.schedules, eventID)
	scheduled.timer.Stop()

	return true
}

// StopAll cancels all timers and stops all pollers
func (m *PollManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for eventID := range m.pollers {
		m.stopLocked(eventID)
	}

	for eventID := range m.schedules {
		m.cancelScheduleLocked(eventID)
	}
}

func (m *PollManager) ActiveCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.pollers)
}

func (m *PollManager) ScheduledCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.schedules)
}

func (m *PollManager) IsRunning(eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.pollers[eventID]
	return exists
}

func (m *PollManager) IsScheduled(eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.schedules[eventID]
	return exists
}
