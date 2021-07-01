package genre

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

// Repository states the required methods from the persistence layer to satisfy business requirements.
type Repository interface {
	// List should returns all Genres if there are no errors.
	List(ctx context.Context) ([]entity.Genre, error)
}

// Driver is an interface described the contract required to satisfy business usecases.
type Driver interface {
	// ListGenres should fetch all Genres and perform intermediary business logic, if any.
	ListGenres(ctx context.Context) ([]entity.Genre, error)
}
