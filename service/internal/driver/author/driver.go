package author

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

func (d *driver) List(ctx context.Context) ([]entity.Author, error) {
	return d.repository.List(ctx)
}
