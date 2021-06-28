package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
)

type authorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) (*authorRepository, error) {
	if db == nil {
		return nil, ErrInvalidDependency
	}

	return &authorRepository{
		db: db,
	}, nil
}

// List selects all Authors in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *authorRepository) List(ctx context.Context) ([]entity.Author, error) {
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

	return authors, nil
}
