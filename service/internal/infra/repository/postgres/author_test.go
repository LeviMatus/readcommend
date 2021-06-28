package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorRepository(t *testing.T) {

	var db sql.DB

	tests := map[string]struct {
		input        *sql.DB
		expect       *authorRepository
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &authorRepository{db: &db},
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewAuthorRepository(tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}

func TestAuthorPostgresRepo_GetAuthors(t *testing.T) {

	var query = "SELECT * FROM author"

	tests := map[string]struct {
		expect               []entity.Author
		setQueryExpectations func(*sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery
		errAssertion         assert.ErrorAssertionFunc
	}{
		"query returns error": {
			expect:       nil,
			errAssertion: assert.Error,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				return query.WillReturnError(errors.New("unable to perform query"))
			},
		},
		"error while scanning result-set": {
			expect:       nil,
			errAssertion: assert.Error,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
					AddRow(42, "john", "doe").RowError(0, errors.New("unable to scan row 0"))
				return query.WillReturnRows(rows)
			},
		},

		"successful get authors": {
			expect: []entity.Author{{
				ID:        42,
				FirstName: "john",
				LastName:  "doe",
			}},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
					AddRow(42, "john", "doe")
				return query.WillReturnRows(rows)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock := newMock(t)
			repo := &authorRepository{db}

			tt.setQueryExpectations(mock.ExpectQuery(regexp.QuoteMeta(query)))

			actual, err := repo.GetAuthors(context.Background())
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}
