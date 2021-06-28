package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/encoding"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
)

type Era struct {
	entity.Era
	MinYear encoding.NullInt16
	MaxYear encoding.NullInt16
}

func (e Era) toEraEntity() entity.Era {
	var (
		min *int16
		max *int16
	)

	if val, _ := e.MinYear.Value(); val == nil {
		min = nil
	} else {
		min = &e.MaxYear.Int16
	}

	if val, _ := e.MaxYear.Value(); val == nil {
		min = nil
	} else {
		max = &e.MaxYear.Int16
	}

	return entity.Era{
		ID:      e.ID,
		MinYear: min,
		MaxYear: max,
	}
}

type eraRepository struct {
	db *sql.DB
}

func NewEraRepository(db *sql.DB) (*eraRepository, error) {
	if db == nil {
		return nil, ErrInvalidDependency
	}

	return &eraRepository{
		db: db,
	}, nil
}

// GetEras selects all Eras in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *eraRepository) GetEras(ctx context.Context) ([]entity.Era, error) {
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

	var eras []entity.Era
	for rows.Next() {
		var era entity.Era
		if err = rows.Scan(&era.ID, &era.Title, &era.MinYear, &era.MaxYear); err != nil {
			return nil, fmt.Errorf("unable to scan data into era: %w", err)
		}
		eras = append(eras, era)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return eras, nil
}
