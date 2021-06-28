package author

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type Repository interface {
	List(ctx context.Context) ([]entity.Author, error)
}

type Driver interface {
	ListAuthors(ctx context.Context)
}
