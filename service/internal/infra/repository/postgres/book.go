package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type bookRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewBookRepository accepts a pointer to a sql.DB type. If the pointer is nil, then an error is returned.
// Otherwise the pointer is wrapped in a bookRepository and a pointer to it is returned.
func NewBookRepository(db *sql.DB, logger *zap.Logger) (*bookRepository, error) {
	if db == nil || logger == nil {
		return nil, ErrInvalidDependency
	}

	return &bookRepository{
		db:     db,
		logger: logger,
	}, nil
}

// Search selects all Books in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *bookRepository) Search(ctx context.Context, params book.SearchInput) ([]entity.Book, error) {
	r.logger.Debug("searching books from postgres repository")

	/*
	 * Start building SQL query
	 */
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
		builder = builder.PlaceholderFormat(sq.Dollar).Where(sq.Eq{"book.title": *params.Title})
	}

	builder = whereInt16Between(builder, "pages", params.MinPages, params.MaxPages)
	builder = whereInt16Between(builder, "year_published", params.MinYearPublished, params.MaxYearPublished)

	if params.Limit != nil {
		builder = builder.Limit(*params.Limit)
	}

	query, values, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SQL query: %w", err)
	}
	r.logger.Debug(fmt.Sprintf("search book query: %s\n search book values: %v\n", query, values))
	/*
	 * Finish building SQL query
	 */

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("unable to get books: %w", err)
	}
	defer rows.Close()

	var books []entity.Book

	// Iterate over result-set, map to entity.Book, and place in resulting slice.
	for rows.Next() {
		var b entity.Book
		if err = rows.Scan(&b.ID,
			&b.Title,
			&b.YearPublished,
			&b.Rating,
			&b.Pages,
			&b.Author.ID,
			&b.Author.FirstName,
			&b.Author.LastName,
			&b.Genre.ID,
			&b.Genre.Title); err != nil {
			return nil, fmt.Errorf("unable to scan data into b: %w", err)
		}
		books = append(books, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	r.logger.Debug(fmt.Sprintf("found %d books in postgres repository", len(books)))

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
