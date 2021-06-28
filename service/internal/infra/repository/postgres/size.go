package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/encoding"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
)

type Size struct {
	ID       int32
	Title    string
	MinPages encoding.NullInt16
	MaxPages encoding.NullInt16
}

func (s Size) toSizeEntity() entity.Size {
	var (
		min *int16
		max *int16
	)

	if val, _ := s.MinPages.Value(); val == nil {
		min = nil
	} else {
		min = &s.MinPages.Int16
	}

	if val, _ := s.MaxPages.Value(); val == nil {
		max = nil
	} else {
		max = &s.MaxPages.Int16
	}

	return entity.Size{
		ID:       s.ID,
		Title:    s.Title,
		MinPages: min,
		MaxPages: max,
	}
}

type sizeRepository struct {
	db *sql.DB
}

func NewSizeRepository(db *sql.DB) (*sizeRepository, error) {
	if db == nil {
		return nil, ErrInvalidDependency
	}

	return &sizeRepository{
		db: db,
	}, nil
}

// GetSizes selects all Sizes in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *sizeRepository) GetSizes(ctx context.Context) ([]entity.Size, error) {
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

	var sizes []entity.Size
	for rows.Next() {
		var size Size
		if err = rows.Scan(&size.ID, &size.Title, &size.MinPages, &size.MaxPages); err != nil {
			return nil, fmt.Errorf("unable to scan data into a genre: %w", err)
		}
		sizes = append(sizes, size.toSizeEntity())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sizes, nil
}
