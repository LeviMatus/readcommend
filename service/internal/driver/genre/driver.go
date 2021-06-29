package genre

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type driver struct {
	repository Repository
}

func NewDriver(r Repository) *driver {
	return &driver{repository: r}
}

func (d *driver) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	return d.repository.List(ctx)
}
