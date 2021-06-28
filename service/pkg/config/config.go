package config

type Config struct {
	Database
	API
}

type Database struct {
	Host         string
	Port         string
	Database     string
	Schema       string
	SSL          string
	Username     string
	Password     string
	DefaultLimit uint64
}

type API struct {
	Port      string
	Interface string
}
