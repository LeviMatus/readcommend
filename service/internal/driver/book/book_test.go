package book_test

import (
	"context"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
)

type inMemoryRepository struct {
	resource map[int32]entity.Book
}

func (r *inMemoryRepository) Search(_ context.Context, _ book.SearchInput) ([]entity.Book, error) {
	var data []entity.Book
	for _, a := range r.resource {
		data = append(data, a)
	}
	return data, nil
}

func TestDriver_Search(t *testing.T) {

	b := entity.Book{
		ID:            1,
		Title:         "The Silmarillion",
		YearPublished: 1977,
		Rating:        3.9,
		Pages:         365,
		Genre: entity.Genre{
			ID:    2,
			Title: "SciFi/Fantasy",
		},
		Author: entity.Author{
			ID:        1,
			FirstName: "John",
			LastName:  "Tolkien",
		},
	}

	repo := inMemoryRepository{resource: map[int32]entity.Book{1: b}}

	driver := book.NewDriver(&repo)
	res, err := driver.SearchBooks(context.Background(), book.SearchInput{})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, b)
}
