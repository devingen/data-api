package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/data-api/dto"
	"net/http"
)

func (controller ServiceController) Update(ctx context.Context, req core.Request) (*core.Response, error) {

	_, webhookStatusCode, webhookError := controller.WebhookService.Pre(ctx, req)
	if webhookError != nil {
		return &core.Response{
			StatusCode: webhookStatusCode,
			Body:       webhookError,
		}, nil
	}

	var body dto.UpdateConfig
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "base-missing-in-path")
	}

	collection, hasCollection := req.PathParameters["collection"]
	if !hasCollection {
		return nil, core.NewError(http.StatusInternalServerError, "collection-missing-in-path")
	}

	id, hasID := req.PathParameters["id"]
	if !hasID {
		return nil, core.NewError(http.StatusBadRequest, "id-missing")
	}

	updated, revision, err := controller.DataService.Update(ctx, base, collection, id, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       dto.UpdateResponse{Updated: updated, Revision: revision},
	}, nil
}
