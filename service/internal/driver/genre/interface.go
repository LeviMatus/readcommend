package genre

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type Repository interface {
	List(ctx context.Context) ([]entity.Genre, error)
}

type Driver interface {
	ListGenres(ctx context.Context) ([]entity.Genre, error)
}
