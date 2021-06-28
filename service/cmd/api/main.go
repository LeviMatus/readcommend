package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/LeviMatus/readcommend/service/internal/api"
	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/infra/repository/postgres"
	"github.com/LeviMatus/readcommend/service/pkg/config"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.Config{
		Database: config.Database{
			Host:     "localhost",
			Port:     "5432",
			Database: "readcommend",
			Schema:   "public",
			SSL:      "disable",
			Username: "postgres",
			Password: "password123",
		},
		API: config.API{
			Port:      "5000",
			Interface: "0.0.0.0",
		},
	}

	conStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database, cfg.Database.SSL)

	if cfg.Database.Schema != "" {
		conStr = fmt.Sprintf("%s search_path=%s", conStr, cfg.Database.Schema)
	}

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	bookRepo, err := postgres.NewBookRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	authorRepo, err := postgres.NewAuthorRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	genreRepo, err := postgres.NewGenreRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	eraRepo, err := postgres.NewEraRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	sizeRepo, err := postgres.NewSizeRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	d, err := driver.New(
		author.NewDriver(authorRepo),
		genre.NewDriver(genreRepo),
		size.NewDriver(sizeRepo),
		era.NewDriver(eraRepo),
		book.NewDriver(bookRepo))

	r, err := api.New(d, cfg.API)
	if err != nil {
		log.Fatal(err)
	}

	l, _ := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.API.Interface, cfg.API.Port))
	r.Serve(l)
}
