package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewBookRepository(t *testing.T) {

	var db sql.DB

	tests := map[string]struct {
		input        *sql.DB
		expect       *bookRepository
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &bookRepository{db: &db},
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewBookRepository(tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}

func TestBookPostgresRepo_GetBooks(t *testing.T) {

	var (
		title                = "The Silmarillion"
		maxYear       int16  = 1980
		minYear       int16  = 1970
		maxPages      int16  = 400
		minPages      int16  = 300
		fantasyGenre  int16  = 2
		fictionGenre  int16  = 7
		johnID        int16  = 42
		christopherID int16  = 43
		limit         uint64 = 25

		silmarillion = entity.Book{
			ID:            1000,
			Title:         title,
			YearPublished: 1977,
			Rating:        3.9,
			Pages:         365,
			Genre:         entity.Genre{ID: 2},
			Author:        entity.Author{ID: 43},
		}
	)

	tests := map[string]struct {
		input                book.SearchInput
		expectedQuery        string
		expect               []entity.Book
		setQueryExpectations func(*sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery
		errAssertion         assert.ErrorAssertionFunc
	}{
		"query returns error": {
			errAssertion: assert.Error,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				return query.WillReturnError(errors.New("unable to perform query"))
			},
		},
		"error while scanning result-set": {
			errAssertion: assert.Error,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
					AddRow(42, "john", "doe").RowError(0, errors.New("unable to scan row 0"))
				return query.WillReturnRows(rows)
			},
		},
		"successful get authors with no parameters": {
			expectedQuery: "SELECT book.id, book.title, year_published, rating, pages, author.id, first_name, " +
				"last_name, genre.id, genre.title FROM book LEFT JOIN author ON book.author_id = author.id " +
				"LEFT JOIN genre ON book.genre_id = genre.id ORDER BY rating DESC",
			expect:       []entity.Book{silmarillion},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"book.id", "book.title", "year_published", "rating", "pages", "author.id", "first_name", "last_name", "genre.id", "genre.title"}).
					AddRow(silmarillion.ID, silmarillion.Title, silmarillion.YearPublished, silmarillion.Rating,
						silmarillion.Pages, silmarillion.Author.ID, silmarillion.Author.FirstName, silmarillion.Author.LastName,
						silmarillion.Genre.ID, silmarillion.Genre.Title)
				return query.WillReturnRows(rows)
			},
		},
		"successful get authors with all parameters": {
			input: book.SearchInput{
				Title:            &title,
				MaxYearPublished: &maxYear,
				MinYearPublished: &minYear,
				MaxPages:         &maxPages,
				MinPages:         &minPages,
				GenreIDs:         []int16{fantasyGenre, fictionGenre},
				AuthorIDs:        []int16{johnID, christopherID},
				Limit:            &limit,
			},
			expectedQuery: "SELECT book.id, book.title, year_published, rating, pages, author.id, first_name, " +
				"last_name, genre.id, genre.title " +
				"FROM book LEFT JOIN author ON book.author_id = author.id LEFT JOIN genre ON book.genre_id = genre.id " +
				"WHERE author_id IN ($1,$2) AND genre_id IN ($3,$4) AND book.title = $5 AND pages >= $6 " +
				"AND pages <= $7 AND year_published >= $8 AND year_published <= $9 ORDER BY rating DESC LIMIT 25",
			expect:       []entity.Book{silmarillion},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"book.id", "book.title", "year_published", "rating", "pages", "author.id", "first_name", "last_name", "genre.id", "genre.title"}).
					AddRow(silmarillion.ID, silmarillion.Title, silmarillion.YearPublished, silmarillion.Rating,
						silmarillion.Pages, silmarillion.Author.ID, silmarillion.Author.FirstName, silmarillion.Author.LastName,
						silmarillion.Genre.ID, silmarillion.Genre.Title)
				return query.WillReturnRows(rows)
			},
		},
		"successful get authors with only lower bounds": {
			input: book.SearchInput{
				MinYearPublished: &minYear,
				MinPages:         &minPages,
			},
			expectedQuery: "SELECT book.id, book.title, year_published, rating, pages, author.id, first_name, " +
				"last_name, genre.id, genre.title " +
				"FROM book LEFT JOIN author ON book.author_id = author.id LEFT JOIN genre ON book.genre_id = genre.id " +
				"WHERE pages >= $1 AND year_published >= $2 ORDER BY rating DESC",
			expect:       []entity.Book{silmarillion},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"book.id", "book.title", "year_published", "rating", "pages", "author.id", "first_name", "last_name", "genre.id", "genre.title"}).
					AddRow(silmarillion.ID, silmarillion.Title, silmarillion.YearPublished, silmarillion.Rating,
						silmarillion.Pages, silmarillion.Author.ID, silmarillion.Author.FirstName, silmarillion.Author.LastName,
						silmarillion.Genre.ID, silmarillion.Genre.Title)
				return query.WillReturnRows(rows)
			},
		},
		"successful get authors with only upper bounds": {
			input: book.SearchInput{
				MaxYearPublished: &maxYear,
				MaxPages:         &maxPages,
			},
			expectedQuery: "SELECT book.id, book.title, year_published, rating, pages, author.id, first_name, " +
				"last_name, genre.id, genre.title " +
				"FROM book LEFT JOIN author ON book.author_id = author.id LEFT JOIN genre ON book.genre_id = genre.id " +
				"WHERE pages <= $1 AND year_published <= $2 ORDER BY rating DESC",
			expect:       []entity.Book{silmarillion},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"book.id", "book.title", "year_published", "rating", "pages", "author.id", "first_name", "last_name", "genre.id", "genre.title"}).
					AddRow(silmarillion.ID, silmarillion.Title, silmarillion.YearPublished, silmarillion.Rating,
						silmarillion.Pages, silmarillion.Author.ID, silmarillion.Author.FirstName, silmarillion.Author.LastName,
						silmarillion.Genre.ID, silmarillion.Genre.Title)
				return query.WillReturnRows(rows)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock := newMock(t)
			repo := &bookRepository{db}

			tt.setQueryExpectations(mock.ExpectQuery(regexp.QuoteMeta(tt.expectedQuery)))

			actual, err := repo.Search(context.Background(), tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}
}
