package era

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

// ListEras fetches entity.Era types from the repository and returns them.
func (d *driver) ListEras(ctx context.Context) ([]entity.Era, error) {
	return d.repository.List(ctx)
}
