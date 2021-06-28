package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
)

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) (*bookRepository, error) {
	if db == nil {
		return nil, ErrInvalidDependency
	}

	return &bookRepository{
		db: db,
	}, nil
}

// Search selects all Books in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *bookRepository) Search(ctx context.Context, params book.SearchInput) ([]entity.Book, error) {

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("book.id", "book.title", "year_published", "rating",
			"pages", "author.id", "first_name", "last_name", "genre.id", "genre.title").
		From("book").
		LeftJoin("author ON book.author_id = author.id").
		LeftJoin("genre ON book.genre_id = genre.id").
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

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		if err = rows.Scan(&book.ID,
			&book.Title,
			&book.YearPublished,
			&book.Rating,
			&book.Pages,
			&book.Author.ID,
			&book.Author.FirstName,
			&book.Author.LastName,
			&book.Genre.ID,
			&book.Genre.Title); err != nil {
			return nil, fmt.Errorf("unable to scan data into book: %w", err)
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
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
