package config

// Config defines the configurations needed to run the readcommend backend service.
type Config struct {
	// Database defines the configs needed to connect to the persistance-layer DB.
	Database Database `mapstructure:"database"`

	// API defines the configs needed to stand up an API listening for network connections.
	API API `mapstructure:"api"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Schema   string `mapstructure:"schema"`
	SSL      string `mapstructure:"ssl-mode"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type API struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}
