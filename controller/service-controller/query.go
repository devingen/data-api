package service_controller

import (
	"github.com/devingen/data-api/dto"
)

func (controller ServiceController) Query(request *dto.QueryRequest) (dto.QueryResponse, error) {

	results, meta, err := controller.Service.Query(
		request.Base,
		request.Collection,
		request.QueryConfig,
	)
	return dto.QueryResponse{Results: results, Meta: meta}, err
}
