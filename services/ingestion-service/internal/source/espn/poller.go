package espn

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

type Poller struct {
	client Client

	gameInfo GameInfo

	interval time.Duration
	output   chan<- domain.RawPlay

	lastSequence int64
}

func newPoller(
	client Client,
	game GameInfo,
	interval time.Duration,
	output chan<- domain.RawPlay,
) *Poller {
	return &Poller{
		client:       client,
		gameInfo:     game,
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

}

func (p *Poller) poll(ctx context.Context) error {
	response, err := p.client.GetPlays(
		ctx,
		p.gameInfo.EventID,
		p.gameInfo.CompetitionID,
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

		rawPlay, err := MapPlay(p.gameInfo.EventID, dto)
		if err != nil {
			return fmt.Errorf(
				"map play %q with sequence %d: %w",
				dto.ID,
				sequence,
				err,
			)
		}

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
