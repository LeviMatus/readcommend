package book

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

// Repository states the required methods from the persistence layer to satisfy business requirements.
type Repository interface {
	// Search should accepts SearchInput items and returns a slice of entity.Book types if no error.
	Search(ctx context.Context, params SearchInput) ([]entity.Book, error)
}

// Driver is an interface described the contract required to satisfy business usecases.
type Driver interface {
	// SearchBooks should fetch all entity.Book types and perform intermediary business logic, if any.
	SearchBooks(ctx context.Context, params SearchInput) ([]entity.Book, error)
}
