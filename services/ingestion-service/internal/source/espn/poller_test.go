package espn

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

type mockClient struct {
	playResponse PlayResponse
	playErr      error

	receivedEventID       string
	receivedCompetitionID string
	getPlaysCallCount     int
}

func (m *mockClient) GetScoreboard(
	ctx context.Context,
	date time.Time,
) (ScoreboardResponse, error) {
	panic("Getscoreboard shouldn't be called by poller")
}

func (m *mockClient) GetPlays(
	ctx context.Context,
	eventID string,
	competitionID string,
) (PlayResponse, error) {
	m.receivedEventID = eventID
	m.receivedCompetitionID = competitionID
	m.getPlaysCallCount++

	return m.playResponse, m.playErr
}

func TestPoller_FirstPollPublishAll(t *testing.T) {
	client := &mockClient{
		playResponse: PlayResponse{
			Items: []PlayDTO{
				newTestPlay("play-103", 103),
				newTestPlay("play-101", 101),
				newTestPlay("play-102", 102),
			},
		},
	}

	output := make(chan domain.RawPlay, 3)

	game := GameInfo{
		EventID:       "event-123",
		CompetitionID: "competition-456",
	}

	poller := Poller{
		client:   client,
		game:     game,
		interval: time.Second,
		output:   output,
	}

	err := poller.poll(context.Background())
	if err != nil {
		t.Fatalf("poll() returned unexpected error: %v", err)
	}

	first := <-output
	second := <-output
	third := <-output

	if first.Sequence != 101 {
		t.Errorf(
			"first published sequence = %d; want 101",
			first.Sequence,
		)
	}

	if second.Sequence != 102 {
		t.Errorf(
			"second published sequence = %d; want 102",
			second.Sequence,
		)
	}

	if third.Sequence != 103 {
		t.Errorf(
			"third published sequence = %d; want 103",
			third.Sequence,
		)
	}

	if first.EventID != "play-101" {
		t.Errorf(
			"first published event ID = %q; want %q",
			first.EventID,
			"play-101",
		)
	}

	if first.GameID != game.EventID {
		t.Errorf(
			"published game ID = %q; want %q",
			first.GameID,
			game.EventID,
		)
	}

	if first.TeamID != "10" {
		t.Errorf(
			"published team ID = %q; want %q",
			first.TeamID,
			"10",
		)
	}

	if poller.lastSequence != 103 {
		t.Errorf(
			"lastSequence = %d; want 103",
			poller.lastSequence,
		)
	}

}

func TestPoller_DuplicateSequenceIgnored(t *testing.T) {
	client := &mockClient{
		playResponse: PlayResponse{
			Items: []PlayDTO{
				newTestPlay("play-103", 103),
				newTestPlay("play-101", 101),
				newTestPlay("play-102", 102),
			},
		},
	}

	output := make(chan domain.RawPlay, 3)

	game := GameInfo{
		EventID:       "event-123",
		CompetitionID: "competition-456",
	}

	poller := NewPoller(
		client,
		game,
		time.Second,
		output,
	)

	if err := poller.poll(context.Background()); err != nil {
		t.Fatal(err)
	}

	for len(output) > 0 {
		<-output
	}

	client.playResponse = PlayResponse{
		Items: []PlayDTO{
			newTestPlay("play-101", 101),
			newTestPlay("play-102", 102),
			newTestPlay("play-103", 103),
			newTestPlay("play-104", 104),
			newTestPlay("play-105", 105),
		},
	}

	if err := poller.poll(context.Background()); err != nil {
		t.Fatal(err)
	}

	first := <-output
	second := <-output

	if first.Sequence != 104 {
		t.Fatalf(
			"first sequence = %d, want 104",
			first.Sequence,
		)
	}

	if second.Sequence != 105 {
		t.Fatalf(
			"second sequence = %d, want 105",
			second.Sequence,
		)
	}

	if poller.lastSequence != 105 {
		t.Fatalf(
			"last sequence = %d, want 105",
			poller.lastSequence,
		)
	}
}

func TestPoller_ContextCanceld(t *testing.T) {
	client := &mockClient{
		playResponse: PlayResponse{
			Items: []PlayDTO{
				newTestPlay("play-103", 103),
			},
		},
	}

	// unbuffered channel
	output := make(chan domain.RawPlay)

	poller := NewPoller(
		client,
		GameInfo{
			EventID:       "event",
			CompetitionID: "competition",
		},
		time.Second,
		output,
	)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error)

	go func() {
		done <- poller.poll(ctx)
	}()

	time.Sleep(20 * time.Millisecond)

	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf(
				"error = %v, want context.Canceled",
				err,
			)
		}
	case <-time.After(time.Second):
		t.Fatal("poll() did not exit after context cancellation")
	}
}

func newTestPlay(
	id string,
	sequence int64,
) PlayDTO {
	return PlayDTO{
		ID:             id,
		SequenceNumber: FlexiableInt64(sequence),

		Text:      "Test play",
		ShortText: "Test",

		Clock: ClockDTO{
			Value:        120,
			DisplayValue: "2:00",
		},

		Period: PeriodDTO{
			Number:       1,
			DisplayValue: "1st Quarter",
		},

		Type: PlayTypeDTO{
			ID:   "1",
			Text: "Test Event",
		},

		Valid:       true,
		ScoringPlay: false,
		ScoreValue:  0,

		Modified: "2026-07-26T20:00Z",

		HomeScore: 10,
		AwayScore: 8,

		Team: RefDTO{
			Ref: "https://sports.core.api.espn.com/v2/sports/basketball/leagues/wnba/seasons/2026/teams/10",
		},
	}
}

func TestPoller_RealPlays(t *testing.T) {
	client := NewClient(ClientConfig{
		PlaysLimits: 3,
	})

	output := make(chan domain.RawPlay, 3)

	game := GameInfo{
		EventID:       "401857026",
		CompetitionID: "401857026",
	}

	poller := Poller{
		client:   client,
		game:     game,
		interval: time.Second,
		output:   output,
	}

	if err := poller.poll(context.Background()); err != nil {
		t.Fatalf("poll() returned unexpected error: %v", err)
	}

	first := <-output
	second := <-output
	third := <-output

	t.Logf("first publish play: homescore %d : awayScore %d, clock: %s", first.HomeScore, first.AwayScore, first.Clock)
	t.Logf("second publish play: homescore %d : awayScore %d, clock: %s", second.HomeScore, second.AwayScore, second.Clock)
	t.Logf("third publish play: homescore %d : awayScore %d, clock: %s", third.HomeScore, third.AwayScore, third.Clock)

	t.Logf("first publish play event: %s, seqeunce: %d", first.Description, first.Sequence)
	t.Logf("second publish play event: %s, seqeunce: %d", second.Description, second.Sequence)
	t.Logf("third publish play event: %s, seqeunce: %d", third.Description, third.Sequence)
}
