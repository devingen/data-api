package aws

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/wrapper"
	"github.com/devingen/data-api/config"
	"github.com/devingen/data-api/controller"
	service_controller "github.com/devingen/data-api/controller/service-controller"
	mongo_data_service "github.com/devingen/data-api/data-service/mongo-data-service"
	http_webhook_service "github.com/devingen/data-api/webhook-service/webhook-interceptor-service"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"log"
)

var db *database.Database

// InitDeps creates the dependencies for the AWS Lambda functions.
func InitDeps() (controller.IServiceController, func(f core.Controller) wrapper.AWSLambdaHandler) {

	validate := validator.New()
	core.SetValidator(validate)

	var appConfig config.ApiConfig
	err := envconfig.Process("data_api", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := logrus.New()
	logrusLevel, parseLevelErr := logrus.ParseLevel(appConfig.LogLevel)
	if parseLevelErr == nil {
		logger.SetLevel(logrusLevel)
	}

	dataService := mongo_data_service.New(getDatabase(appConfig))
	interceptorService := http_webhook_service.New(appConfig.Webhook.URL, appConfig.Webhook.Headers)
	serviceController := service_controller.New(appConfig, dataService, interceptorService)

	wrap := generateWrapper(appConfig)

	return serviceController, wrap
}

func getDatabase(appConfig config.ApiConfig) *database.Database {
	if db == nil {
		var err error
		db, err = database.New(appConfig.Mongo.URI)
		if err != nil {
			log.Fatalf("Database connection failed when creating a new database %s", err.Error())
		}
	} else if !db.IsConnected() {
		err := db.ConnectWithURI(appConfig.Mongo.URI)
		if err != nil {
			log.Fatalf("Database connection failed when connecting to an existing database %s", err.Error())
		}
	}
	return db
}

func generateWrapper(appConfig config.ApiConfig) func(f core.Controller) wrapper.AWSLambdaHandler {
	return func(f core.Controller) wrapper.AWSLambdaHandler {
		ctx := context.Background()

		// add logger
		withLogger := wrapper.WithLogger(appConfig.LogLevel, f)

		// convert to HTTP handler
		handler := wrapper.WithLambdaHandler(ctx, withLogger)
		return handler
	}
}
