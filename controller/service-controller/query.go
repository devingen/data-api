package service_controller

import (
	"github.com/devingen/data-api/dto"
	"github.com/devingen/data-api/hooks"
)

func (controller ServiceController) Query(request *dto.QueryRequest) (dto.QueryResponse, error) {

	err := hooks.CheckEligibility(request.AuthorizationHeader)
	if err != nil {
		return dto.QueryResponse{}, err
	}

	results, meta, err := controller.Service.Query(
		request.Base,
		request.Collection,
		request.QueryConfig,
	)
	return dto.QueryResponse{Results: results, Meta: meta}, err
}
