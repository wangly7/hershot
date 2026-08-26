package espn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
	"github.com/wangly7/hershot/shared/events"
)

type Poller struct {
	client Client

	game GameInfo

	streamID   string
	streamMode events.StreamMode

	interval time.Duration
	output   chan<- domain.RawPlay

	lastSequence int64
}

func NewPoller(
	client Client,
	game GameInfo,
	streamID string,
	streamMode events.StreamMode,
	interval time.Duration,
	output chan<- domain.RawPlay,
) *Poller {
	return &Poller{
		client:       client,
		game:         game,
		streamID:     streamID,
		streamMode:   streamMode,
		interval:     interval,
		output:       output,
		lastSequence: -1,
	}
}

// Run starts polling immediately and then continue polling at the
// configured interval until the context is canceled.
//
// A temporary ESPN or mapping error does not stop the Poller. It is logged,
// and the poller tries again on the next tick.
func (p *Poller) Run(ctx context.Context) {
	// immidately starts polling
	if err := p.poll(ctx); err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		log.Printf(
			"poll ESPN game eventID=%s, competitionID=%s: %v",
			p.game.EventID,
			p.game.CompetitionID,
			err,
		)
	}

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				if errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) {
					return
				}
				log.Printf(
					"poll ESPN game eventID=%s competition=%s: %v",
					p.game.EventID,
					p.game.CompetitionID,
					err,
				)
			}
		}
	}
}

func (p *Poller) poll(ctx context.Context) error {
	response, err := p.client.GetPlays(
		ctx,
		p.game.EventID,
		p.game.CompetitionID,
	)
	if err != nil {
		return err
	}

	// Do not assume that ESPN always returns items in ascending order
	sort.SliceStable(response.Items, func(i, j int) bool {
		return int64(response.Items[i].SequenceNumber) < int64(response.Items[j].SequenceNumber)
	})

	for _, dto := range response.Items {
		sequence := dto.SequenceNumber

		if sequence <= FlexiableInt64(p.lastSequence) {
			continue
		}

		rawPlay, err := MapPlay(p.game.EventID, dto)

		if err != nil {
			return fmt.Errorf(
				"map play %q with sequence %d: %w",
				dto.ID,
				sequence,
				err,
			)
		}

		rawPlay.StreamID = p.streamID
		rawPlay.StreamMode = p.streamMode

		// Use select instead of a plain:
		//
		// p.out <- RawPlay
		// Otherwise the poller could remain blocked forever when the output channle
		// has no receiver, even after ctx is canceled.
		select {
		case p.output <- rawPlay:
			p.lastSequence = int64(sequence)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
