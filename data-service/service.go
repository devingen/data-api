package data_service

import (
	"context"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/data-api/dto"
)

type IDataService interface {
	Query(ctx context.Context, base, collection string, config *coremodel.QueryConfig) ([]*coremodel.DataModel, *coremodel.Meta, error)
	Update(ctx context.Context, base, collection, id string, config *dto.UpdateConfig) (string, int, error)
	Create(ctx context.Context, base, collection string, config *dto.CreateConfig) (string, int, error)
	Delete(ctx context.Context, base, collection, id string) (string, int, error)
}
