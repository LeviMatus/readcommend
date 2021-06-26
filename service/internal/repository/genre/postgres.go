package genre

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type genrePostgresRepo struct {
	db *sql.DB
}

// Compile-time check to ensure genrePostgresRepo satisfies the Repository interface.
var _ Repository = (*genrePostgresRepo)(nil)

// NewPostgresRepository accepts a Ptr to a sql.DB. If the Ptr is nil, an error will be thrown.
// The returned repository interfaces with Postgres as its DB resource.
func NewPostgresRepository(db *sql.DB) (*genrePostgresRepo, error) {
	if db == nil {
		return nil, errors.New("expected a non-nil db")
	}

	return &genrePostgresRepo{db: db}, nil
}

// GetGenres selects all Genres in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *genrePostgresRepo) GetGenres(ctx context.Context) ([]Genre, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM genre")
	if err != nil {
		return nil, fmt.Errorf("unable to get genres: %w", err)
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var genre Genre
		if err = rows.Scan(&genre.ID, &genre.Title); err != nil {
			return nil, fmt.Errorf("unable to scan data into a genre: %w", err)
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

// Close terminates the wrapped sql.DB.
func (r *genrePostgresRepo) Close() error {
	return r.db.Close()
}
