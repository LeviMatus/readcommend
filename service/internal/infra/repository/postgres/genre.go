package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
)

type genreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) (*genreRepository, error) {
	if db == nil {
		return nil, ErrInvalidDependency
	}

	return &genreRepository{
		db: db,
	}, nil
}

// List selects all Genres in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *genreRepository) List(ctx context.Context) ([]entity.Genre, error) {
	query, _, err := sq.StatementBuilder.
		Select("*").
		From("genre").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SQL query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to get genres: %w", err)
	}
	defer rows.Close()

	var genres []entity.Genre
	for rows.Next() {
		var genre entity.Genre
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
