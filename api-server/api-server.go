package api_server

import (
	"context"
	"github.com/devingen/data-api/config"
	service_controller "github.com/devingen/data-api/controller/service-controller"
	mongo_data_service "github.com/devingen/data-api/data-service/mongo-data-service"
	http_webhook_service "github.com/devingen/data-api/webhook-service/webhook-interceptor-service"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/wrapper"
	"github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
)

// New creates a new HTTP server
func New(appConfig config.ApiConfig, db *database.Database) *http.Server {

	logger := logrus.New()
	logrusLevel, parseLevelErr := logrus.ParseLevel(appConfig.LogLevel)
	if parseLevelErr == nil {
		logger.SetLevel(logrusLevel)
	}

	validate := validator.New()
	core.SetValidator(validate)

	srv := &http.Server{Addr: ":" + appConfig.Port}

	dataService := mongo_data_service.New(db)
	interceptorService := http_webhook_service.New(appConfig.Webhook.URL, appConfig.Webhook.Headers)
	serviceController := service_controller.New(appConfig, dataService, interceptorService)

	wrap := generateWrapper(appConfig)

	router := mux.NewRouter()
	router.HandleFunc("/{base}/{collection}/query", wrap(serviceController.Query)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/{collection}/create", wrap(serviceController.Create)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/{collection}/{id}/update", wrap(serviceController.Update)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/{collection}/{id}/delete", wrap(serviceController.Delete)).Methods(http.MethodPost)

	http.Handle("/", &server.CORSRouterDecorator{
		R: router,
		Headers: map[string]string{
			server.CORSAccessControlAllowHeaders: server.CORSAccessControlAllowHeadersDefaultValue + ",devingen-product-id" + ",api-key",
			server.CORSAccessControlAllowMethods: server.CORSAccessControlAllowMethodsDefaultValue,
		},
		AllowSenderOrigin: true,
	})
	return srv
}

func generateWrapper(appConfig config.ApiConfig) func(f core.Controller) func(http.ResponseWriter, *http.Request) {
	return func(f core.Controller) func(http.ResponseWriter, *http.Request) {
		ctx := context.Background()

		// add logger
		withLogger := wrapper.WithLogger(appConfig.LogLevel, f)

		// convert to HTTP handler
		handler := wrapper.WithHTTPHandler(ctx, withLogger)
		return handler
	}
}
