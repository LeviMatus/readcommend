package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

const (
	minimumPageParam int16 = 1
	maximumPageParam int16 = 10000
	minimumYearParam int16 = 1800
	maximumYearParam int16 = 2100

	bookSearchParamKey = "book-search-params"
)

func bookRoutes(h *bookHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(
			cors.Handler(cors.Options{AllowedMethods: []string{"GET"}}),
			ValidateBookRequest,
		)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			_ = render.Render(w, r, ErrMethodNotAllowed(r.Method))
			return
		})
		r.Get("/", h.List)
	})
	return r
}

type BookRequest struct {
	_ struct{}

	Title            *string `schema:"title"`
	MaxYearPublished *int16  `schema:"max-year"`
	MinYearPublished *int16  `schema:"min-year"`
	MaxPages         *int16  `schema:"max-pages"`
	MinPages         *int16  `schema:"min-pages"`
	GenreIDs         []int16 `schema:"genres"`
	AuthorIDs        []int16 `schema:"authors"`
	Limit            *uint64 `schema:"limit"`
}

// ValidateBookRequest maps the query parameters to a BookRequest struct, which is injected
// into the context of the request. As a part of this process, BookRequest.GenreIDs and
// BookRequest.AuthorIDs are validated. If a string, such as "alpha" appears in the incoming string list,
// then then validation fails and a 400 StatusCode code is returned.
//
// Following this, the resulting BookRequest is validated. If any search criteria fail validation, then
// the routine returns a 400 StatusCode code and error message.
func ValidateBookRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			_ = render.Render(w, r, ErrBadRequest(fmt.Errorf("an unexpected error occurred: %w", err)))
			return
		}

		queryParams := new(BookRequest)
		if err := schema.NewDecoder().Decode(queryParams, r.Form); err != nil {
			var schemaErr schema.MultiError
			if !errors.As(err, &schemaErr) {
				_ = render.Render(w, r, ErrBadRequest(fmt.Errorf("an unexpected error occured: %w", schemaErr)))
				return
			}
			for k, v := range schemaErr {
				switch v.(type) {
				case schema.ConversionError:
					_ = render.Render(w, r, ErrBadRequest(fmt.Errorf("%w: received wrong type for parameter %s", entity.ErrInvalidQueryParam, k)))
					return
				case schema.UnknownKeyError:
					log.Println("WARN: " + err.Error())
				}
			}
		}

		if queryParams.MinPages != nil && !util.Int16InRange(*queryParams.MinPages, minimumPageParam, maximumPageParam) {
			_ = render.Render(w, r, ErrBadRequest(
				fmt.Errorf("%w: min-pages is %d but should be in range [%d,%d]",
					entity.ErrInvalidQueryParam,
					*queryParams.MinPages,
					minimumPageParam,
					maximumPageParam)))
			return
		}

		if queryParams.MaxPages != nil && !util.Int16InRange(*queryParams.MaxPages, minimumPageParam, maximumPageParam) {
			_ = render.Render(w, r, ErrBadRequest(
				fmt.Errorf("%w: max-pages is %d but should be in range [%d,%d]",
					entity.ErrInvalidQueryParam,
					*queryParams.MaxPages,
					minimumPageParam,
					maximumPageParam)))
			return
		}

		if queryParams.MinYearPublished != nil && !util.Int16InRange(*queryParams.MinYearPublished, minimumYearParam, maximumYearParam) {
			_ = render.Render(w, r, ErrBadRequest(
				fmt.Errorf("%w: min-year is %d but should be in range [%d,%d]",
					entity.ErrInvalidQueryParam,
					*queryParams.MinYearPublished,
					minimumYearParam,
					maximumYearParam)))
			return
		}

		if queryParams.MaxYearPublished != nil && !util.Int16InRange(*queryParams.MaxYearPublished, minimumYearParam, maximumYearParam) {
			_ = render.Render(w, r, ErrBadRequest(
				fmt.Errorf("%w: max-year is %d but should be in range [%d,%d]",
					entity.ErrInvalidQueryParam,
					*queryParams.MaxYearPublished,
					minimumYearParam,
					maximumYearParam)))
			return
		}

		if queryParams.Limit != nil && *queryParams.Limit < 1 {
			_ = render.Render(w, r, ErrBadRequest(
				fmt.Errorf("%w: limit is %d but should be greater than 0",
					entity.ErrInvalidQueryParam,
					*queryParams.Limit)))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), bookSearchParamKey, queryParams)))
	})
}

// bookHandler...
type bookHandler struct {
	driver book.Driver
}

func NewBookHandler(driver book.Driver) (*bookHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil book driver is required to create a book handler")
	}

	return &bookHandler{driver: driver}, nil
}

func (b *bookHandler) List(w http.ResponseWriter, r *http.Request) {
	reqParams, ok := r.Context().Value(bookSearchParamKey).(*BookRequest)

	// This should have been placed into the context by the GET api/v1/books middleware
	if !ok || reqParams == nil {
		_ = render.Render(w, r, ErrInternalServer(errors.New("expected middleware to inject params into context")))
		return
	}

	books, err := b.driver.SearchBooks(r.Context(), book.SearchInput{
		Title:            reqParams.Title,
		MaxYearPublished: reqParams.MaxYearPublished,
		MinYearPublished: reqParams.MinYearPublished,
		MaxPages:         reqParams.MaxPages,
		MinPages:         reqParams.MinPages,
		GenreIDs:         reqParams.GenreIDs,
		AuthorIDs:        reqParams.AuthorIDs,
		Limit:            reqParams.Limit,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalServer(err))
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
