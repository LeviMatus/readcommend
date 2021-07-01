package author

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

// ListAuthors fetches entity.Author types from the repository and returns them.
func (d *driver) ListAuthors(ctx context.Context) ([]entity.Author, error) {
	return d.repository.List(ctx)
}
