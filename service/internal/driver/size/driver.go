package size

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

func (d *driver) ListSizes(ctx context.Context) ([]entity.Size, error) {
	return d.repository.List(ctx)
}
