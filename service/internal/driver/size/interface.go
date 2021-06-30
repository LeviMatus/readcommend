package size

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

// Repository states the required methods from the persistence layer to satisfy business requirements.
type Repository interface {
	// List should return all Sizes.
	List(ctx context.Context) ([]entity.Size, error)
}

// Driver is an interface described the contract required to satisfy business usecases.
type Driver interface {
	// ListSizes should fetch all Sizes an perform intermediary business logic, if any.
	ListSizes(ctx context.Context) ([]entity.Size, error)
}
