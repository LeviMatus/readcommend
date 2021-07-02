package cmd

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/LeviMatus/readcommend/service/internal/api"
	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/infra/repository/postgres"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)

	attachDatabaseFlags(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the readcommend server",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: setup logger
		conStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database, cfg.Database.SSL)

		if cfg.Database.Schema != "" {
			conStr = fmt.Sprintf("%s search_path=%s", conStr, cfg.Database.Schema)
		}

		db, err := sql.Open("postgres", conStr)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		if err := db.Ping(); err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		bookRepo, err := postgres.NewBookRepository(db)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		authorRepo, err := postgres.NewAuthorRepository(db)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		genreRepo, err := postgres.NewGenreRepository(db)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		eraRepo, err := postgres.NewEraRepository(db)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		sizeRepo, err := postgres.NewSizeRepository(db)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		d, err := driver.New(
			author.NewDriver(authorRepo),
			genre.NewDriver(genreRepo),
			size.NewDriver(sizeRepo),
			era.NewDriver(eraRepo),
			book.NewDriver(bookRepo))

		r, err := api.New(d)
		if err != nil {
			// TODO: logging
			ExitRequirements.Exit()
		}

		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port))
		if err != nil {
			// TODO: logging
			ExitListen.Exit()
		}

		if err := r.Serve(l); err != nil {
			// TODO: logging
			ExitServing.Exit()
		}
	},
}
