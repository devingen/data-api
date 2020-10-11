package main

import (
	"github.com/devingen/api-core/database"
	cm "github.com/devingen/api-core/server"
	"github.com/devingen/data-api/controller"
	"github.com/devingen/data-api/server/handler"
	"github.com/devingen/data-api/service"
	"log"
	"net/http"
)

// Runs the server that contains all the services
func main() {

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	databaseService := service.NewDatabaseService(db)
	serviceController := controller.NewServiceController(databaseService)

	// create a Service Handler that uses Database AtamaService
	h := handler.NewHttpServiceHandler(serviceController)

	http.Handle("/", &cm.CORSRouterDecorator{R: h.Router})
	err = http.ListenAndServe(":1003", &cm.CORSRouterDecorator{R: h.Router})
	if err != nil {
		log.Fatalf("Listen and serve failed %s", err.Error())
	}
}
