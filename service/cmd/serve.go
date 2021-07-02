package cmd

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/LeviMatus/readcommend/service/internal/api"
	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/infra/repository/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serveCmd)

	attachDatabaseFlags(serveCmd)

	serveCmd.Flags().StringVar(&cfg.Database.Host,
		"api-host",
		"0.0.0.0",
		`The host that the API listens on (default "0.0.0.0")`)
	serveCmd.Flags().StringVar(&cfg.Database.Port,
		"api-port",
		"5000",
		`The port that the API listens on (default "5000")`)

	viper.BindPFlag("api.host", serveCmd.Flag("api-host"))
	viper.BindPFlag("api.port", serveCmd.Flag("api-port"))
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the readcommend server",
	Run: func(cmd *cobra.Command, args []string) {
		defer logger.Sync()

		conStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database, cfg.Database.SSL)

		if cfg.Database.Schema != "" {
			conStr = fmt.Sprintf("%s search_path=%s", conStr, cfg.Database.Schema)
		}

		db, err := sql.Open("postgres", conStr)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to connect to database: %s", err))
			ExitRequirements.Exit()
		}

		if err := db.Ping(); err != nil {
			logger.Error(fmt.Sprintf("unable to verify DB connection: %s", err))
			ExitRequirements.Exit()
		}
		logger.Info("database connection established")

		bookRepo, err := postgres.NewBookRepository(db, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to create Book repository: %s", err))
			ExitRequirements.Exit()
		}

		authorRepo, err := postgres.NewAuthorRepository(db, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to create Author repository: %s", err))
			ExitRequirements.Exit()
		}

		genreRepo, err := postgres.NewGenreRepository(db, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to create Genre repository: %s", err))
			ExitRequirements.Exit()
		}

		eraRepo, err := postgres.NewEraRepository(db, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to create Era repository: %s", err))
			ExitRequirements.Exit()
		}

		sizeRepo, err := postgres.NewSizeRepository(db, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to create Size repository: %s", err))
			ExitRequirements.Exit()
		}

		r, err := api.New(
			author.NewDriver(authorRepo),
			size.NewDriver(sizeRepo),
			genre.NewDriver(genreRepo),
			era.NewDriver(eraRepo),
			book.NewDriver(bookRepo),
			logger)

		if err != nil {
			logger.Error(fmt.Sprintf("unable to create Driver: %s", err))
			ExitRequirements.Exit()
		}

		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port))
		if err != nil {
			logger.Error(fmt.Sprintf("unable to listen on specified interface/port: %s", err))
			ExitListen.Exit()
		}

		if err := r.Serve(l); err != nil {
			logger.Error(fmt.Sprintf("an error occurred while serving: %s", err))
			ExitServing.Exit()
		}
	},
}
