package era

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

// Repository states the required methods from the persistence layer to satisfy business requirements.
type Repository interface {
	// List should return all Eras.
	List(ctx context.Context) ([]entity.Era, error)
}

// Driver is an interface described the contract required to satisfy business usecases.
type Driver interface {
	// ListEras should fetch all Eras and perform intermediary business logic, if any.
	ListEras(ctx context.Context) ([]entity.Era, error)
}
