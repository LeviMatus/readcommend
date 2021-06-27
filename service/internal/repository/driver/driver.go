package driver

import (
	"database/sql"
	"fmt"

	"github.com/LeviMatus/readcommend/service/pkg/config"
	_ "github.com/lib/pq"
)

type Driver struct {
	postgres *sql.DB
}

func New(config config.Database) (*Driver, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.Username, config.Password, config.Database, config.SSL)

	if config.Schema != "" {
		connStr = fmt.Sprintf("%s search_path=%s", connStr, config.Schema)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to open connection to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to validate database connection: %w", err)
	}

	return &Driver{postgres: db}, nil
}

func (d *Driver) Driver() *sql.DB {
	return d.postgres
}
