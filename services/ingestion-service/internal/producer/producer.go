package producer

import (
	"context"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

type Producer interface {
	Publish(ctx context.Context, play domain.RawPlay) error
	Close()
}
