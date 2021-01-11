package service_controller

import (
	"github.com/devingen/data-api/dto"
	"github.com/devingen/data-api/hooks"
)

func (controller ServiceController) Update(request *dto.UpdateRequest) (dto.UpdateResponse, error) {

	err := hooks.CheckEligibility(request.AuthorizationHeader)
	if err != nil {
		return dto.UpdateResponse{}, err
	}

	updated, revision, err := controller.Service.Update(
		request.Base,
		request.Collection,
		request.ID,
		request.UpdateConfig,
	)
	return dto.UpdateResponse{Updated: updated, Revision: revision}, err
}
