package service_controller

import (
	"github.com/devingen/data-api/config"
	"github.com/devingen/data-api/controller"
	"github.com/devingen/data-api/data-service"
	webhook_service "github.com/devingen/data-api/webhook-service"
)

// ServiceController implements DataController interface by using AtamaService
type ServiceController struct {
	DataService    data_service.IDataService
	WebhookService webhook_service.IDataWebhookService
	EnableCreate   bool
	EnableUpdate   bool
	EnableDelete   bool
}

// New generates new ServiceController
func New(c config.ApiConfig, dataService data_service.IDataService, webhookService webhook_service.IDataWebhookService) controller.IServiceController {
	return ServiceController{
		DataService:    dataService,
		WebhookService: webhookService,
		EnableCreate:   c.EnableCreate,
		EnableUpdate:   c.EnableUpdate,
		EnableDelete:   c.EnableDelete,
	}
}
