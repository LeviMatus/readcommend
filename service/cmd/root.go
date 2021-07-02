package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func init() {
	cobra.OnInitialize(setupConfig)
}

var rootCmd = &cobra.Command{
	Use:   "readcommend",
	Short: "Backend tooling for human bookworms",
	Long: `Tooling and supporting backend API for interfacing with the readcommend backing database.
				use this to query for books directly or standup an application server.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		// for simplicity, I'm just using a default zap production logger.
		logger, err = zap.NewProduction()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "cannot setup logger: %s", err)
			ExitRequirements.Exit()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func setupConfig() {
	// Get the config provided as-is and marshal it.
	// Depending on the inputs provided by the user, they will be over written.
	// Precedence is Flag > Env Var > Config File > Zero Values.
	b, err := yaml.Marshal(cfg)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ExitConfigSetup.Exit()
	}

	configReader := bytes.NewReader(b)
	viper.SetConfigType("yaml")

	// Take the default configReader and populate the viper config we're building with defaults.
	if err := viper.MergeConfig(configReader); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ExitConfigSetup.Exit()
	}

	// if a config file is specified, then use it explicitly. Otherwise check the user's home directory for
	// a file with the default name.
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			ExitConfigSetup.Exit()
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".readcommend")
	}

	// tell viper to merge the values found in the config file into
	// the ongoing config. If no file is found, that's potentially fine. There's still
	// a chance that the runtime will find all needed variables defined via env vars or flags.
	if err := viper.MergeInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			fmt.Printf("user's config file %s was not found, but will try continuing anyway...", configFile)
		default:
			_, _ = fmt.Fprintln(os.Stderr, err)
			ExitConfigSetup.Exit()
		}
	}

	// tell viper to get env vars based on tags on the config struct and move them into the
	// ongoing config. Replace any hyphens or "." with "_" so that our env variables have normal naming
	// conventions.
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_", "-", "_"),
	)

	// Finally unmarshal viper's config into the application config type.
	if err := viper.Unmarshal(&cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ExitConfigSetup.Exit()
	}

	// Lastly, as a special case, if the user want's to define a password for the database
	// in the CLI interactively, do so here. This step occurs last because it is a CLI argument and
	// should overwrite anything defined in the environment or config file.
	if promptDatabasePass {
		validatePassword()
	}
}
