package service

import (
	coremodel "github.com/devingen/api-core/model"
)

type DataService interface {
	Query(base, collection string, config *coremodel.QueryConfig) ([]*coremodel.DataModel, *coremodel.Meta, error)
}
