package config

type Config struct {
	Database
	API
}

type Database struct {
	Host     string
	Port     string
	Database string
	Schema   string
	SSL      string
	Username string
	Password string
}

type API struct {
	Port         string
	Interface    string
	DefaultLimit uint64
	Timeout      int16
}
