package era_test

import (
	"context"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
)

type inMemoryRepository struct {
	resource map[int32]entity.Era
}

func (r *inMemoryRepository) List(_ context.Context) ([]entity.Era, error) {
	var data []entity.Era
	for _, e := range r.resource {
		data = append(data, e)
	}
	return data, nil
}

func TestDriver_List(t *testing.T) {

	e := entity.Era{ID: 1, Title: "Modern"}

	repo := inMemoryRepository{resource: map[int32]entity.Era{1: e}}

	driver := era.NewDriver(&repo)
	res, err := driver.ListEras(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, e)
}
