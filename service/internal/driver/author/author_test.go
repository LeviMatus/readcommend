package author_test

import (
	"context"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
)

type inMemoryRepository struct {
	resource map[int32]entity.Author
}

func (r *inMemoryRepository) List(_ context.Context) ([]entity.Author, error) {
	var data []entity.Author
	for _, a := range r.resource {
		data = append(data, a)
	}
	return data, nil
}

func TestDriver_List(t *testing.T) {

	a := entity.Author{ID: 1, FirstName: "John", LastName: "Tolkien"}

	repo := inMemoryRepository{resource: map[int32]entity.Author{1: a}}

	driver := author.NewDriver(&repo)
	res, err := driver.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, a)
}
