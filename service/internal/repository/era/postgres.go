package era

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type eraPostgresRepo struct {
	db *sql.DB
}

// Compile-time check to ensure eraPostgresRepo satisfies the Repository interface.
var _ Repository = (*eraPostgresRepo)(nil)

// NewPostgresRepository accepts a Ptr to a sql.DB. If the Ptr is nil, an error will be thrown.
// The returned repository interfaces with Postgres as its DB resource.
func NewPostgresRepository(db *sql.DB) (*eraPostgresRepo, error) {
	if db == nil {
		return nil, errors.New("expected a non-nil db")
	}

	return &eraPostgresRepo{db: db}, nil
}

// GetEras selects all Eras in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *eraPostgresRepo) GetEras(ctx context.Context) ([]Era, error) {
	query, _, err := sq.StatementBuilder.
		Select("*").
		From("era").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SQL query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to get eras: %w", err)
	}
	defer rows.Close()

	var eras []Era
	for rows.Next() {
		var era Era
		if err = rows.Scan(&era.ID, &era.Title, &era.StartYear, &era.EndYear); err != nil {
			return nil, fmt.Errorf("unable to scan data into era: %w", err)
		}
		eras = append(eras, era)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return eras, nil
}

// Close terminates the wrapped sql.DB.
func (r *eraPostgresRepo) Close() error {
	return r.db.Close()
}
