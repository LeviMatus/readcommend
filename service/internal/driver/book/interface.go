package book

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type Repository interface {
	Search(ctx context.Context, params SearchInput) ([]entity.Book, error)
}

type Driver interface {
	SearchBooks(ctx context.Context, params SearchInput)
}
