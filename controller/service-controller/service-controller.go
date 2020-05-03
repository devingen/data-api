package service_controller

import (
	"github.com/devingen/data-api/service"
)

// ServiceController implements DataController interface by using AtamaService
type ServiceController struct {
	Service service.DataService
}
