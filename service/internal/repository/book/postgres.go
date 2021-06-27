package book

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type bookPostgresRepo struct {
	db *sql.DB
}

// Compile-time check to ensure bookPostgresRepo satisfies the Repository interface.
var _ Repository = (*bookPostgresRepo)(nil)

// NewPostgresRepository accepts a Ptr to a sql.DB. If the Ptr is nil, an error will be thrown.
// The returned repository interfaces with Postgres as its DB resource.
func NewPostgresRepository(db *sql.DB) (*bookPostgresRepo, error) {
	if db == nil {
		return nil, errors.New("expected a non-nil db")
	}

	return &bookPostgresRepo{db: db}, nil
}

// GetBooks selects all Books in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *bookPostgresRepo) GetBooks(ctx context.Context, params GetBooksParams) ([]Book, error) {

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").
		From("book").
		OrderBy("rating DESC")

	builder = whereInt16In(builder, "author_id", params.AuthorIDs)
	builder = whereInt16In(builder, "genre_id", params.GenreIDs)

	if params.Title != nil {
		builder = builder.PlaceholderFormat(sq.Dollar).Where(sq.Eq{"title": *params.Title})
	}

	if params.Rating != nil {
		builder = builder.PlaceholderFormat(sq.Dollar).Where(sq.Eq{"rating": *params.Rating})
	}

	builder = whereInt16Between(builder, "pages", params.MinPages, params.MaxPages)
	builder = whereInt16Between(builder, "year_published", params.MinYearPublished, params.MaxYearPublished)

	var limit uint64 = 10
	if params.Limit != nil {
		limit = *params.Limit
	}
	builder = builder.Limit(limit)

	query, values, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SQL query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("unable to get books: %w", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err = rows.Scan(&book.ID,
			&book.Title,
			&book.YearPublished,
			&book.Rating,
			&book.Pages,
			&book.GenreID,
			&book.AuthorID); err != nil {
			return nil, fmt.Errorf("unable to scan data into book: %w", err)
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// Close terminates the wrapped sql.DB.
func (r *bookPostgresRepo) Close() error {
	return r.db.Close()
}

// whereInt16In accepts a query builder, a target column, and a slice of int16s.
// If the slice is non-empty, then it a SQL WHERE clause section will be added for
// records with col values IN the provided ints slice. The mutated builder is returned.
func whereInt16In(builder sq.SelectBuilder, col string, ints []int16) sq.SelectBuilder {
	if len(ints) > 0 {
		var values = make([]interface{}, len(ints))
		var placeholders = make([]interface{}, len(ints))
		for i := range values {
			values[i] = ints[i]
			placeholders[i] = "?"
		}
		return builder.
			PlaceholderFormat(sq.Dollar).
			Where(
				fmt.Sprintf("%s IN (%s)", col, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(placeholders)), ","), "[]")),
				values...,
			)
	}

	return builder
}

// whereInt16Between accepts a query builder, a target column, and pointers to int16
// for the min and max values of the range. Because the builder dependency does not support
// Postgres' BETWEEN semantics, this is done as two WHERE clause filters. Its possible
// to let min or max be nil. In such a scenario, the query would only provide a lower
// or upper bound.
func whereInt16Between(builder sq.SelectBuilder, col string, min *int16, max *int16) sq.SelectBuilder {
	if min != nil {
		builder = builder.PlaceholderFormat(sq.Dollar).Where(sq.GtOrEq{col: *min})
	}

	if max != nil {
		builder = builder.PlaceholderFormat(sq.Dollar).Where(sq.LtOrEq{col: *max})
	}

	return builder
}
