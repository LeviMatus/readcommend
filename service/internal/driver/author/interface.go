package author

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

// Repository states the required methods from the persistence layer to satisfy business requirements.
type Repository interface {
	// List should return all Authors if there are no errors.
	List(ctx context.Context) ([]entity.Author, error)
}

// Driver is an interface described the contract required to satisfy business usecases.
type Driver interface {
	// ListAuthors should fetch all Authors and performs intermediary business logic, if any.
	ListAuthors(ctx context.Context) ([]entity.Author, error)
}
