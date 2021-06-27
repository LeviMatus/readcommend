package size

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type sizePostgresRepo struct {
	db *sql.DB
}

// Compile-time check to ensure sizePostgresRepo satisfies the Repository interface.
var _ Repository = (*sizePostgresRepo)(nil)

// NewPostgresRepository accepts a Ptr to a sql.DB. If the Ptr is nil, an error will be thrown.
// The returned repository interfaces with Postgres as its DB resource.
func NewPostgresRepository(db *sql.DB) (*sizePostgresRepo, error) {
	if db == nil {
		return nil, errors.New("expected a non-nil db")
	}

	return &sizePostgresRepo{db: db}, nil
}

// GetSizes selects all Sizes in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *sizePostgresRepo) GetSizes(ctx context.Context) ([]Size, error) {
	query, _, err := sq.StatementBuilder.
		Select("*").
		From("size").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SQL query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to get sizes: %w", err)
	}
	defer rows.Close()

	var sizes []Size
	for rows.Next() {
		var size Size
		if err = rows.Scan(&size.ID, &size.Title, &size.MinimumPages, &size.MaximumPages); err != nil {
			return nil, fmt.Errorf("unable to scan data into a genre: %w", err)
		}
		sizes = append(sizes, size)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sizes, nil
}

// Close terminates the wrapped sql.DB.
func (r *sizePostgresRepo) Close() error {
	return r.db.Close()
}
