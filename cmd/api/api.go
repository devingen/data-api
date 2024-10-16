package main

import (
	"github.com/devingen/api-core/database"
	"github.com/devingen/data-api/api-server"
	"github.com/devingen/data-api/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logrus.New()

	var appConfig config.ApiConfig
	err := envconfig.Process("data_api", &appConfig)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(appConfig.Mongo.URI)
	logger.Info("Running server on port ", appConfig.Port)

	db, err := database.New(appConfig.Mongo.URI)
	if err != nil {
		logger.Fatalf("Database connection failed %s", err.Error())
	}
	logger.Info("Connected to database")

	srv := api_server.New(appConfig, db)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("Listen and serve failed %s", err)
	}
}
