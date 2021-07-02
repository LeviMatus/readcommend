package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	promptDatabasePass bool
)

// attachDatabaseFlags can be used by more commands in the future. It attaches all database param flags
// to a specified command.
func attachDatabaseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Database.Host,
		"db-host",
		"localhost",
		`The host where the backend DB lives (default "localhost")`)
	cmd.Flags().StringVar(&cfg.Database.Port,
		"db-port",
		"5432",
		`The port which the bcakend DB listens on (default "5432")`)
	cmd.Flags().StringVar(&cfg.Database.Database,
		"db-name",
		"readcommend",
		`The name of the database to connect to (default "readcommend").`)
	cmd.Flags().StringVar(&cfg.Database.Schema,
		"db-schema",
		"public",
		`The database schema to use (default "public")`)
	cmd.Flags().StringVar(&cfg.Database.SSL,
		"db-ssl-mode",
		"disable",
		`Whether to use SSL connection to the database (default "disable")`)
	cmd.Flags().StringVar(&cfg.Database.Username,
		"db-username",
		"postgres",
		`The username to use when connecting to the backend DB (default "postgres")`)
	cmd.Flags().BoolVar(&promptDatabasePass,
		"db-password",
		false,
		`Prompt for the database password, if true`)

	viper.BindPFlag("database.host", cmd.Flag("db-host"))
	viper.BindPFlag("database.port", cmd.Flag("db-port"))
	viper.BindPFlag("database.name", cmd.Flag("db-name"))
	viper.BindPFlag("database.schema", cmd.Flag("db-schema"))
	viper.BindPFlag("database.ssl-mode", cmd.Flag("db-ssl-mode"))
	viper.BindPFlag("database.username", cmd.Flag("db-username"))
}

// validatePassword will prompt for a user-provided password interactively.
// Password input will be hidden from the terminal.
func validatePassword() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Database Password",
		Validate: validate,
		Mask:     ' ',
	}

	result, err := prompt.Run()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ExitConfigSetup.Exit()
	}

	cfg.Database.Password = result
}
