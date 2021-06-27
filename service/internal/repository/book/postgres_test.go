package book

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to create sql mock: %v", err)
	}
	return db, mock
}

func TestNewPostgresRepository(t *testing.T) {

	var db sql.DB

	tests := map[string]struct {
		input        *sql.DB
		expect       *bookPostgresRepo
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &bookPostgresRepo{db: &db},
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewPostgresRepository(tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}

func TestBookPostgresRepo_GetBooks(t *testing.T) {

	var (
		title                 = "The Silmarillion"
		maxYear       int16   = 1980
		minYear       int16   = 1970
		rating        float32 = 3.9
		maxPages      int16   = 400
		minPages      int16   = 300
		fantasyGenre  int16   = 2
		fictionGenre  int16   = 7
		johnID        int16   = 42
		christopherID int16   = 43
		limit         uint64  = 25

		silmarillion = Book{
			ID:            1000,
			Title:         title,
			YearPublished: 1977,
			Rating:        3.9,
			Pages:         365,
			GenreID:       2,
			AuthorID:      christopherID,
		}
	)

	tests := map[string]struct {
		input                GetBooksParams
		expectedQuery        string
		expect               []Book
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
			expectedQuery: "SELECT * FROM book ORDER BY rating DESC LIMIT 10",
			expect:        []Book{silmarillion},
			errAssertion:  assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "title", "year_published", "rating", "pages", "genre_id", "author_id"}).
					AddRow(silmarillion.ID, silmarillion.Title, silmarillion.YearPublished, silmarillion.Rating,
						silmarillion.Pages, silmarillion.GenreID, silmarillion.AuthorID)
				return query.WillReturnRows(rows)
			},
		},
		"successful get authors with all parameters": {
			input: GetBooksParams{
				Title:            &title,
				MaxYearPublished: &maxYear,
				MinYearPublished: &minYear,
				MaxPages:         &maxPages,
				MinPages:         &minPages,
				Rating:           &rating,
				GenreIDs:         []int16{fantasyGenre, fictionGenre},
				AuthorIDs:        []int16{johnID, christopherID},
				Limit:            &limit,
			},
			expectedQuery: "SELECT * FROM book WHERE author_id IN ($1,$2) AND genre_id IN ($3,$4) AND title = $5 AND " +
				"rating = $6 AND pages >= $7 AND pages <= $8 AND year_published >= $9 AND year_published <= $10 " +
				"ORDER BY rating DESC LIMIT 25",
			expect:       []Book{silmarillion},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "title", "year_published", "rating", "pages", "genre_id", "author_id"}).
					AddRow(silmarillion.ID, silmarillion.Title, silmarillion.YearPublished, silmarillion.Rating,
						silmarillion.Pages, silmarillion.GenreID, silmarillion.AuthorID)
				return query.WillReturnRows(rows)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock := newMock(t)
			repo := &bookPostgresRepo{db}
			defer func() {
				repo.Close()
			}()

			tt.setQueryExpectations(mock.ExpectQuery(regexp.QuoteMeta(tt.expectedQuery)))

			actual, err := repo.GetBooks(context.Background(), tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}
}
