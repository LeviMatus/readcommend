package size

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

// ListSizes fetches entity.Size types from the repository and returns them.
func (d *driver) ListSizes(ctx context.Context) ([]entity.Size, error) {
	return d.repository.List(ctx)
}
