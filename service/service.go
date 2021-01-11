package service

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/data-api/dto"
)

type DataService interface {
	Query(base, collection string, config *coremodel.QueryConfig) ([]*coremodel.DataModel, *coremodel.Meta, error)
	Update(base, collection, id string, config *dto.UpdateConfig) (string, int, error)
}
