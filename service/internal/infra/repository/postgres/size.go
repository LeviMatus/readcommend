package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/encoding"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

// size is a persistence layer model. It has support for nullable SQL fields.
type size struct {
	ID       int32
	Title    string
	MinPages encoding.NullInt16
	MaxPages encoding.NullInt16
}

func (s size) toSizeEntity() entity.Size {
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
	db     *sql.DB
	logger *zap.Logger
}

// NewSizeRepository accepts a pointer to a sql.DB type. If the pointer is nil, then an error is returned.
// Otherwise the pointer is wrapped in an sizeRepository and a pointer to it is returned.
func NewSizeRepository(db *sql.DB, logger *zap.Logger) (*sizeRepository, error) {
	if db == nil || logger == nil {
		return nil, ErrInvalidDependency
	}

	return &sizeRepository{
		db:     db,
		logger: logger,
	}, nil
}

// List selects all Sizes in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *sizeRepository) List(ctx context.Context) ([]entity.Size, error) {
	r.logger.Debug("listing sizes from postgres repository")

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

	// Iterate over result-set, map to entity.Size, and place in resulting slice.
	for rows.Next() {
		var size size
		if err = rows.Scan(&size.ID, &size.Title, &size.MinPages, &size.MaxPages); err != nil {
			return nil, fmt.Errorf("unable to scan data into a genre: %w", err)
		}
		sizes = append(sizes, size.toSizeEntity())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	r.logger.Debug(fmt.Sprintf("found %d sizes in postgres repository", len(sizes)))

	return sizes, nil
}
