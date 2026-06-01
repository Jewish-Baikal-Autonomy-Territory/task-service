package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Searcher interface {
	Search(ctx context.Context, filter *Filter) ([]uuid.UUID, time.Time, error)
}
