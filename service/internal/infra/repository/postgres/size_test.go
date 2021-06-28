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

func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to create sql mock: %v", err)
	}
	return db, mock
}

func TestNewSizeRepository(t *testing.T) {

	var db sql.DB

	tests := map[string]struct {
		input        *sql.DB
		expect       *sizeRepository
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &sizeRepository{db: &db},
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewSizeRepository(tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}

func TestSizeRepository_GetSizes(t *testing.T) {

	var query = "SELECT * FROM size"

	tests := map[string]struct {
		expect               []entity.Size
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

		"successful get sizes": {
			expect: []entity.Size{
				{
					ID:    0,
					Title: "Any",
				}, {
					ID:       1,
					Title:    "Short story – up to 35 pages",
					MaxPages: util.Int16Ptr(34),
				}, {
					ID:       2,
					Title:    "Novelette – 35 to 85 pages",
					MinPages: util.Int16Ptr(35),
					MaxPages: util.Int16Ptr(84),
				}, {
					ID:       3,
					Title:    "Novella – 85 to 200 pages",
					MinPages: util.Int16Ptr(85),
					MaxPages: util.Int16Ptr(199),
				}, {
					ID:       4,
					Title:    "Novel – 200 to 500 pages",
					MinPages: util.Int16Ptr(200),
					MaxPages: util.Int16Ptr(499),
				}, {
					ID:       5,
					Title:    "Brick – 500 to 800 pages",
					MinPages: util.Int16Ptr(500),
					MaxPages: util.Int16Ptr(799),
				}, {
					ID:       6,
					Title:    "Monument – 800 pages and up",
					MinPages: util.Int16Ptr(800),
				},
			},
			errAssertion: assert.NoError,
			setQueryExpectations: func(query *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				rows := sqlmock.NewRows([]string{"id", "title", "min_pages", "max_pages"}).
					AddRow(0, "Any", nil, nil).
					AddRow(1, "Short story – up to 35 pages", nil, 34).
					AddRow(2, "Novelette – 35 to 85 pages", 35, 84).
					AddRow(3, "Novella – 85 to 200 pages", 85, 199).
					AddRow(4, "Novel – 200 to 500 pages", 200, 499).
					AddRow(5, "Brick – 500 to 800 pages", 500, 799).
					AddRow(6, "Monument – 800 pages and up", 800, nil)
				return query.WillReturnRows(rows)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock := newMock(t)
			repo := &sizeRepository{db}

			tt.setQueryExpectations(mock.ExpectQuery(regexp.QuoteMeta(query)))

			actual, err := repo.GetSizes(context.Background())
			assert.True(t, reflect.DeepEqual(tt.expect, actual))
			tt.errAssertion(t, err)
		})
	}
}
