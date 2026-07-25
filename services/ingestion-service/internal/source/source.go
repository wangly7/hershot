package source

import (
	"context"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

type Source interface {
	Run(
		ctx context.Context,
		output chan<- domain.RawPlay,
	) error
}
