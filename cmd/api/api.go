package main

import (
	"log"
	"net/http"

	api_server "github.com/devingen/data-api/api-server"
	"github.com/devingen/data-api/config"

	"github.com/devingen/api-core/database"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	var appConfig config.ApiConfig
	err := envconfig.Process("data_api", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Running server on port", appConfig.Port)

	db, err := database.New(appConfig.Mongo.URI)
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	srv := api_server.New(appConfig, db)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and serve failed %s", err.Error())
	}
}
