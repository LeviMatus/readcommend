package main

import (
	"log"

	"github.com/LeviMatus/readcommend/service/internal/api"
	"github.com/LeviMatus/readcommend/service/internal/repository/driver"
	"github.com/LeviMatus/readcommend/service/pkg/config"
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

	dbDriver, err := driver.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	r, err := api.New(dbDriver.Driver(), cfg.API)
	if err != nil {
		log.Fatal(err)
	}

	r.Listen()
}
