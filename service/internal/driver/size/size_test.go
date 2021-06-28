package size_test

import (
	"context"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
)

type inMemoryRepository struct {
	resource map[int32]entity.Size
}

func (r *inMemoryRepository) List(_ context.Context) ([]entity.Size, error) {
	var data []entity.Size
	for _, s := range r.resource {
		data = append(data, s)
	}
	return data, nil
}

func TestDriver_List(t *testing.T) {

	s := entity.Size{ID: 1, Title: "Any"}

	repo := inMemoryRepository{resource: map[int32]entity.Size{1: s}}

	driver := size.NewDriver(&repo)
	res, err := driver.ListSizes(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, s)
}
