package v1

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/LeviMatus/readcommend/service/internal/api/params"
	"github.com/LeviMatus/readcommend/service/internal/domain"
	"github.com/LeviMatus/readcommend/service/internal/repository/book"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// Book...
type Book struct {
	repo book.Repository
}

func NewBookHandler(db *sql.DB) (*Book, error) {
	if db == nil {
		return nil, errors.New("non-nil database client is required to create a book handler")
	}
	repo, err := book.NewPostgresRepository(db)
	if err != nil {
		return nil, fmt.Errorf("unable to create a book handler: %w", err)
	}

	return &Book{repo: repo}, nil
}

func (b *Book) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

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

	title = params.String(ctx, "title")

	minPages, err := params.Int16(ctx, "min_pages")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	maxPages, err = params.Int16(ctx, "max_pages")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	maxYear, err = params.Int16(ctx, "max_year")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	minYear, err = params.Int16(ctx, "min_year")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	authorIDs, err = params.Int16Slice(r, "author_ids")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	genreIDs, err = params.Int16Slice(r, "genre_ids")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	limit, err = params.Uint64(ctx, "limit")
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	books, err := b.repo.GetBooks(r.Context(), book.GetBooksParams{
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
	var out = make([]domain.Book, len(books))
	for i, b := range books {
		out[i] = domain.Book{
			ID:            b.ID,
			Title:         b.Title,
			YearPublished: b.YearPublished,
			Rating:        b.Rating,
			Pages:         b.Pages,
			Genre: domain.Genre{
				ID:    b.Genre.ID,
				Title: b.Genre.Title,
			},
			Author: domain.Author{
				ID:        b.Author.ID,
				FirstName: b.Author.FirstName,
				LastName:  b.Author.LastName,
			},
		}
	}

	render.JSON(w, r, out)
}
