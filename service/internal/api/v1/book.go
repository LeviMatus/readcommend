package v1

import (
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/api/params"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// Book...
type Book struct {
	driver book.Driver
}

func NewBookHandler(driver book.Driver) (*Book, error) {
	if driver == nil {
		return nil, errors.New("non-nil book driver is required to create a book handler")
	}

	return &Book{driver: driver}, nil
}

func (b *Book) List(w http.ResponseWriter, r *http.Request) {
	var (
		title     *string
		minPages  *int16
		maxPages  *int16
		maxYear   *int16
		minYear   *int16
		genreIDs  []int16
		authorIDs []int16
		limit     *uint64
	)

	title = params.String(r, "title")

	minPages, err := params.Int16(r, "min_pages")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	maxPages, err = params.Int16(r, "max_pages")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	maxYear, err = params.Int16(r, "max_year")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	minYear, err = params.Int16(r, "min_year")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	authorIDs, err = params.Int16Slice(r, "authors")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	genreIDs, err = params.Int16Slice(r, "genres")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	limit, err = params.Uint64(r, "limit")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	books, err := b.driver.SearchBooks(r.Context(), book.SearchInput{
		Title:            title,
		MaxYearPublished: maxYear,
		MinYearPublished: minYear,
		MaxPages:         maxPages,
		MinPages:         minPages,
		GenreIDs:         genreIDs,
		AuthorIDs:        authorIDs,
		Limit:            limit,
	})
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// An adaptor between the service layer and persistance layer
	// wouldn't be out of the question, but the conversion is very simple
	// so I'll just do it directly here. In the future, abstracting this
	// may be appropriate.
	var out = make([]entity.Book, len(books))
	for i, b := range books {
		out[i] = entity.Book{
			ID:            b.ID,
			Title:         b.Title,
			YearPublished: b.YearPublished,
			Rating:        b.Rating,
			Pages:         b.Pages,
			Genre: entity.Genre{
				ID:    b.Genre.ID,
				Title: b.Genre.Title,
			},
			Author: entity.Author{
				ID:        b.Author.ID,
				FirstName: b.Author.FirstName,
				LastName:  b.Author.LastName,
			},
		}
	}

	render.JSON(w, r, out)
}
