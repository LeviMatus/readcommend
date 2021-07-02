package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type authorRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAuthorRepository accepts a pointer to a sql.DB type. If the pointer is nil, then an error is returned.
// Otherwise the pointer is wrapped in an authorRepository and a pointer to it is returned.
func NewAuthorRepository(db *sql.DB, logger *zap.Logger) (*authorRepository, error) {
	if db == nil || logger == nil {
		return nil, ErrInvalidDependency
	}

	return &authorRepository{
		db:     db,
		logger: logger,
	}, nil
}

// List selects all Authors in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *authorRepository) List(ctx context.Context) ([]entity.Author, error) {
	r.logger.Debug("listing authors from postgres repository")
	query, _, err := sq.StatementBuilder.
		Select("*").
		From("author").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SQL query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to get authors: %w", err)
	}
	defer rows.Close()

	var authors []entity.Author

	// Iterate over result-set, map to entity.Author, and place in resulting slice.
	for rows.Next() {
		var author entity.Author
		if err = rows.Scan(&author.ID, &author.FirstName, &author.LastName); err != nil {
			return nil, fmt.Errorf("unable to scan data into author: %w", err)
		}
		authors = append(authors, author)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	r.logger.Debug(fmt.Sprintf("found %d authors in postgres repository", len(authors)))

	return authors, nil
}
