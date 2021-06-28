package genre_test

import (
	"context"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
)

type inMemoryRepository struct {
	resource map[int32]entity.Genre
}

func (r *inMemoryRepository) List(_ context.Context) ([]entity.Genre, error) {
	var data []entity.Genre
	for _, g := range r.resource {
		data = append(data, g)
	}
	return data, nil
}

func TestDriver_List(t *testing.T) {

	g := entity.Genre{ID: 1, Title: "SciFi/Fantasy"}

	repo := inMemoryRepository{resource: map[int32]entity.Genre{1: g}}

	driver := genre.NewDriver(&repo)
	res, err := driver.ListGenres(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, g)
}
