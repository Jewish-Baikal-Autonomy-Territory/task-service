package task

import (
	"context"

	"github.com/google/uuid"
)

type Searcher interface {
	Search(ctx context.Context, filter *Filter) ([]uuid.UUID, error)
}
