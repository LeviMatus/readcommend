package postgres

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewEraRepository(t *testing.T) {

	var db sql.DB

	tests := map[string]struct {
		input        *sql.DB
		expect       *eraRepository
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &eraRepository{db: &db},
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewEraRepository(tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}

func TestEraPostgresRepo_GetEras(t *testing.T) {

	var query = "SELECT * FROM era"

	tests := map[string]struct {
		expect               []entity.Era
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
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(42, "SciFi/Fantasy").RowError(0, errors.New("unable to scan row 0"))
				return query.WillReturnRows(rows)
			},
		},

		"successful get eras": {
			expect: []entity.Era{
				{
					ID:    0,
					Title: "Any",
				}, {
					ID:      1,
					Title:   "Classic",
					MaxYear: util.Int16Ptr(1969),
				}, {
					ID:      2,
					Title:   "Modern",
					MinYear: util.Int16Ptr(1970),
				},
			},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "title", "min_year", "max_year"}).
					AddRow(0, "Any", nil, nil).
					AddRow(1, "Classic", nil, 1969).
					AddRow(2, "Modern", 1970, nil)
				return query.WillReturnRows(rows)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock := newMock(t)
			repo := &eraRepository{db}

			tt.setQueryExpectations(mock.ExpectQuery(regexp.QuoteMeta(query)))

			actual, err := repo.List(context.Background())
			assert.True(t, reflect.DeepEqual(tt.expect, actual))
			tt.errAssertion(t, err)
		})
	}
}
