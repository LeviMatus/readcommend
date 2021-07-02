package config

// Config defines the configurations needed to run the readcommend backend service.
type Config struct {
	// Database defines the configs needed to connect to the persistance-layer DB.
	Database Database `mapstructure:"database"`

	// API defines the configs needed to stand up an API listening for network connections.
	API API `mapstructure:"api"`
}

type Database struct {
	// Host the database listens on.
	Host string `mapstructure:"host"`
	// Port the database listens on.
	Port string `mapstructure:"port"`
	// Database name.
	Database string `mapstructure:"database"`
	// Schema to use in the database.
	Schema string `mapstructure:"schema"`
	// SSL mode.
	SSL string `mapstructure:"ssl-mode"`
	// Username to connect with.
	Username string `mapstructure:"username"`
	// Password to connect with.
	Password string `mapstructure:"password"`
}

type API struct {
	// Port the API will listen on.
	Port string `mapstructure:"port"`
	// Host the API will listen on.
	Host string `mapstructure:"host"`
}
