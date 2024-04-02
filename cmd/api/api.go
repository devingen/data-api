package main

import (
	api_server "github.com/devingen/data-api/api-server"
	"github.com/devingen/data-api/config"
	"log"
	"net/http"

	"github.com/devingen/api-core/database"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	var appConfig config.ApiConfig
	err := envconfig.Process("data_api", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.New(appConfig.Mongo.URI)
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	srv := api_server.New(appConfig, db)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and serve failed %s", err.Error())
	}
}
