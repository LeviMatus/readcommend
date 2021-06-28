package size

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type Repository interface {
	List(ctx context.Context) ([]entity.Size, error)
}

type Driver interface {
	ListSizes(ctx context.Context)
}
