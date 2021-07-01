package api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	books = []entity.Book{
		{
			ID:            1,
			Title:         "The Hobbit",
			YearPublished: 1937,
			Rating:        4.3,
			Pages:         310,
			Genre:         entity.Genre{ID: 2, Title: "Fantasy/SciFy"},
			Author: entity.Author{
				ID:        1,
				FirstName: "John",
				LastName:  "Tolkien",
			},
		},
		{
			ID:            2,
			Title:         "The Lord of the Rings",
			YearPublished: 1954,
			Rating:        4.5,
			Pages:         1178,
			Genre:         entity.Genre{ID: 2, Title: "Fantasy"},
			Author: entity.Author{
				ID:        1,
				FirstName: "John",
				LastName:  "Tolkien",
			},
		},
	}
	authors = []entity.Author{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Tolkien",
		}, {
			ID:        2,
			FirstName: "Christopher",
			LastName:  "Tolkien",
		},
	}
	genres = []entity.Genre{
		{ID: 1, Title: "Young Adult"},
		{ID: 2, Title: "Fantasy/SciFy"},
		{ID: 3, Title: "Romance"},
		{ID: 4, Title: "Nonfiction"},
		{ID: 5, Title: "Mystery"},
		{ID: 6, Title: "Memoir"},
		{ID: 7, Title: "Fiction"},
		{ID: 8, Title: "Childrens"},
	}
	sizes = []entity.Size{
		{ID: 0, Title: "Any"},
		{ID: 1, Title: "Short Story", MaxPages: util.Int16Ptr(34)},
		{ID: 2, Title: "Novelette", MinPages: util.Int16Ptr(35), MaxPages: util.Int16Ptr(84)},
		{ID: 3, Title: "Novella", MinPages: util.Int16Ptr(85), MaxPages: util.Int16Ptr(199)},
		{ID: 4, Title: "Novel", MinPages: util.Int16Ptr(200), MaxPages: util.Int16Ptr(499)},
		{ID: 5, Title: "Brick", MinPages: util.Int16Ptr(500), MaxPages: util.Int16Ptr(799)},
		{ID: 6, Title: "Monument", MinPages: util.Int16Ptr(800)},
	}
	eras = []entity.Era{
		{ID: 1, Title: "Any"},
		{ID: 2, Title: "Classic", MaxYear: util.Int16Ptr(1969)},
		{ID: 3, Title: "Modern", MinYear: util.Int16Ptr(1970)},
	}
)

func BenchmarkAPI_Books(b *testing.B) {
	driver := drivertest.DriverMock{}
	driver.
		On("SearchBooks", mock.MatchedBy(func(_ context.Context) bool { return true }), book.SearchInput{}).
		Return(books, nil)

	apiServer, err := New(&driver)
	assert.NoError(b, err)
	assert.NotNil(b, apiServer)

	testServer := httptest.NewServer(apiServer.mux)
	defer testServer.Close()

	for i := 0; i < b.N; i++ {
		_, _ = testServer.Client().Get(fmt.Sprintf("%s/api/v1/books", testServer.URL))
	}
}

func BenchmarkAPI_Authors(b *testing.B) {
	driver := drivertest.DriverMock{}
	driver.
		On("ListAuthors", mock.MatchedBy(func(_ context.Context) bool { return true })).
		Return(authors, nil)

	apiServer, err := New(&driver)
	assert.NoError(b, err)
	assert.NotNil(b, apiServer)

	testServer := httptest.NewServer(apiServer.mux)
	defer testServer.Close()

	for i := 0; i < b.N; i++ {
		_, _ = testServer.Client().Get(fmt.Sprintf("%s/api/v1/authors", testServer.URL))
	}
}

func BenchmarkAPI_Genres(b *testing.B) {
	driver := drivertest.DriverMock{}
	driver.
		On("ListGenres", mock.MatchedBy(func(_ context.Context) bool { return true })).
		Return(genres, nil)

	apiServer, err := New(&driver)
	assert.NoError(b, err)
	assert.NotNil(b, apiServer)

	testServer := httptest.NewServer(apiServer.mux)
	defer testServer.Close()

	for i := 0; i < b.N; i++ {
		_, _ = testServer.Client().Get(fmt.Sprintf("%s/api/v1/genres", testServer.URL))
	}
}

func BenchmarkAPI_Sizes(b *testing.B) {
	driver := drivertest.DriverMock{}
	driver.
		On("ListSizes", mock.MatchedBy(func(_ context.Context) bool { return true })).
		Return(sizes, nil)

	apiServer, err := New(&driver)
	assert.NoError(b, err)
	assert.NotNil(b, apiServer)

	testServer := httptest.NewServer(apiServer.mux)
	defer testServer.Close()

	for i := 0; i < b.N; i++ {
		_, _ = testServer.Client().Get(fmt.Sprintf("%s/api/v1/sizes", testServer.URL))
	}
}

func BenchmarkAPI_Eras(b *testing.B) {
	driver := drivertest.DriverMock{}
	driver.
		On("ListEras", mock.MatchedBy(func(_ context.Context) bool { return true })).
		Return(eras, nil)

	apiServer, err := New(&driver)
	assert.NoError(b, err)
	assert.NotNil(b, apiServer)

	testServer := httptest.NewServer(apiServer.mux)
	defer testServer.Close()

	for i := 0; i < b.N; i++ {
		_, _ = testServer.Client().Get(fmt.Sprintf("%s/api/v1/eras", testServer.URL))
	}
}
