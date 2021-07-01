package genre

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type driver struct {
	repository Repository
}

// NewDriver creates a driver which wraps the repository. The wrapper
// will perform business logic against the usecases of Size entity.
func NewDriver(r Repository) *driver {
	return &driver{repository: r}
}

// ListGenres fetches entity.Genre types from the repository and returns them.
func (d *driver) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	return d.repository.List(ctx)
}
