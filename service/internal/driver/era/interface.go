package era

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type Repository interface {
	List(ctx context.Context) ([]entity.Era, error)
}

type Driver interface {
	ListEras(ctx context.Context) ([]entity.Era, error)
}
