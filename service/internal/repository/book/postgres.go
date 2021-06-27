package book

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
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

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").
		From("book").
		OrderBy("rating DESC")

	if len(params.AuthorIDs) > 0 {
		placeholders := make([]string, len(params.AuthorIDs))
		values := make([]interface{}, len(params.AuthorIDs))
		for i := range params.AuthorIDs {
			placeholders[i] = "?"
			values[i] = params.AuthorIDs[i]
		}
		builder = builder.
			PlaceholderFormat(squirrel.Dollar).
			Where(
				fmt.Sprintf("author_id IN (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(placeholders)), ","), "[]")),
				values...,
			)
	}

	if len(params.GenreIDs) > 0 {
		placeholders := make([]string, len(params.GenreIDs))
		values := make([]interface{}, len(params.GenreIDs))
		for i := range params.GenreIDs {
			placeholders[i] = "?"
			values[i] = params.GenreIDs[i]
		}
		builder = builder.
			PlaceholderFormat(squirrel.Dollar).
			Where(
				fmt.Sprintf("genre_id IN (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(placeholders)), ","), "[]")),
				values...)
	}

	if params.Title != nil {
		builder = builder.PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"title": *params.Title})
	}

	if params.Rating != nil {
		builder = builder.PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"rating": *params.Rating})
	}

	if params.MaxPages != nil {
		builder = builder.PlaceholderFormat(squirrel.Dollar).Where(squirrel.LtOrEq{"pages": *params.MaxPages})
	}

	if params.MinPages != nil {
		builder = builder.PlaceholderFormat(squirrel.Dollar).Where(squirrel.GtOrEq{"pages": *params.MinPages})
	}

	if params.MaxYearPublished != nil {
		builder = builder.PlaceholderFormat(squirrel.Dollar).Where(squirrel.LtOrEq{"year_published": *params.MaxYearPublished})
	}

	if params.MinYearPublished != nil {
		builder = builder.PlaceholderFormat(squirrel.Dollar).Where(squirrel.LtOrEq{"year_published": *params.MinYearPublished})
	}

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
