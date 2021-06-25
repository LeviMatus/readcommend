package genre

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
		expect       *genrePostgresRepo
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &genrePostgresRepo{db: &db},
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

func TestAuthorPostgresRepo_GetAuthors(t *testing.T) {

	var query = "SELECT * FROM genre"

	tests := map[string]struct {
		input                Genre
		expect               []Genre
		setQueryExpectations func(*sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery
		errAssertion         assert.ErrorAssertionFunc
	}{
		"query returns error": {
			input: Genre{
				ID:    42,
				Title: "SciFi/Fantasy",
			},
			expect:       nil,
			errAssertion: assert.Error,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				return query.WillReturnError(errors.New("unable to perform query"))
			},
		},
		"error while scanning result-set": {
			input: Genre{
				ID:    42,
				Title: "SciFi/Fantasy",
			},
			expect:       nil,
			errAssertion: assert.Error,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(42, "SciFi/Fantasy").RowError(0, errors.New("unable to scan row 0"))
				return query.WillReturnRows(rows)
			},
		},

		"successful get authors": {
			input: Genre{
				ID:    42,
				Title: "SciFi/Fantasy",
			},
			expect: []Genre{{
				ID:    42,
				Title: "SciFi/Fantasy",
			}},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(42, "SciFi/Fantasy")
				return query.WillReturnRows(rows)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock := newMock(t)
			repo := &genrePostgresRepo{db}
			defer func() {
				repo.Close()
			}()

			tt.setQueryExpectations(mock.ExpectQuery(regexp.QuoteMeta(query)))

			actual, err := repo.GetGenres(context.Background())
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}
}
