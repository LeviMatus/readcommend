package cmd

import (
	"os"

	"github.com/LeviMatus/readcommend/service/pkg/config"
)

// ExitCode refers to the code returned when the CLI terminates
// its runtime. int8 is used to comply with POSIX standards.
type ExitCode int8

const (
	// OK takes on the value of 0, which is exiting without error.
	// Normally I'd leave the zero-value at a default like "Undefined"
	// to avoid confusion, but in this case the zero-value has meaning.
	OK ExitCode = iota

	// ExitConfigSetup indicates a fatal error took place while setting up the
	// config for executing the CLI and any sub-commands.
	ExitConfigSetup

	// ExitRequirements indicates that a fatal error took place while setting up
	// the requirements of a command, such as establishing drivers and/or repositories.
	ExitRequirements

	// ExitListen indicates that a fatal error took place while establishing a listener
	// against a host/port binding.
	ExitListen

	// ExitServing indicates that a fatal error took place while serving.
	ExitServing
)

// Exit calls the appropriate exit code on the ExitCode type.
func (e ExitCode) Exit() {
	os.Exit(int(e))
}

var (
	cfg        config.Config
	configFile string
)
