package author

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type authorPostgresRepo struct {
	db *sql.DB
}

// Compile-time check to ensure authorPostgresRepo satisfies the Repository interface.
var _ Repository = (*authorPostgresRepo)(nil)

// NewPostgresRepository accepts a Ptr to a sql.DB. If the Ptr is nil, an error will be thrown.
// The returned repository interfaces with Postgres as its DB resource.
func NewPostgresRepository(db *sql.DB) (*authorPostgresRepo, error) {
	if db == nil {
		return nil, errors.New("expected a non-nil db")
	}

	return &authorPostgresRepo{db: db}, nil
}

// GetAuthors selects all Authors in the repository. If the query fails or encounters an erro while
// cursing through the result set, then an error is returned.
func (r *authorPostgresRepo) GetAuthors(ctx context.Context) ([]Author, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM author")
	if err != nil {
		return nil, fmt.Errorf("unable to get authors: %w", err)
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		if err = rows.Scan(&author.ID, &author.FirstName, &author.LastName); err != nil {
			return nil, fmt.Errorf("unable to scan data into author: %w", err)
		}
		authors = append(authors, author)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

// Close terminates the wrapped sql.DB.
func (r *authorPostgresRepo) Close() error {
	return r.db.Close()
}
