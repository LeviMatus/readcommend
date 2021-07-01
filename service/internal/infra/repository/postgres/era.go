package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/encoding"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	sq "github.com/Masterminds/squirrel"
)

// era is a persistence layer model. It has support for nullable SQL fields.
type era struct {
	// ID is the primary identifier of an era.
	ID int32

	// Title is the name of an era.
	Title string

	// MinYear satisfies interfaces necessary to scan SQL's NULL into a Go type.
	MinYear encoding.NullInt16

	// MaxYear satisfies interfaces necessary to scan SQL's NULL into a Go type.
	MaxYear encoding.NullInt16
}

func (e era) toEraEntity() entity.Era {
	var (
		min *int16
		max *int16
	)

	if val, _ := e.MinYear.Value(); val == nil {
		min = nil
	} else {
		min = &e.MinYear.Int16
	}

	if val, _ := e.MaxYear.Value(); val == nil {
		max = nil
	} else {
		max = &e.MaxYear.Int16
	}

	return entity.Era{
		ID:      e.ID,
		Title:   e.Title,
		MinYear: min,
		MaxYear: max,
	}
}

type eraRepository struct {
	db *sql.DB
}

// NewEraRepository accepts a pointer to a sql.DB type. If the pointer is nil, then an error is returned.
// Otherwise the pointer is wrapped in an eraRepository and a pointer to it is returned.
func NewEraRepository(db *sql.DB) (*eraRepository, error) {
	if db == nil {
		return nil, ErrInvalidDependency
	}

	return &eraRepository{
		db: db,
	}, nil
}

// List selects all Eras in the repository. If the query fails or encounters an error while
// cursing through the result set, then an error is returned.
func (r *eraRepository) List(ctx context.Context) ([]entity.Era, error) {
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

	// Iterate over result-set, map to entity.Era, and place in resulting slice.
	for rows.Next() {
		var era era
		if err = rows.Scan(&era.ID, &era.Title, &era.MinYear, &era.MaxYear); err != nil {
			return nil, fmt.Errorf("unable to scan data into era: %w", err)
		}
		eras = append(eras, era.toEraEntity())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return eras, nil
}
