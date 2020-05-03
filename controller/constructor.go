package controller

import (
	"github.com/devingen/data-api/controller/service-controller"
	"github.com/devingen/data-api/service"
)

// NewServiceController generates new ServiceController
func NewServiceController(service service.DataService) *service_controller.ServiceController {
	return &service_controller.ServiceController{
		Service: service,
	}
}
