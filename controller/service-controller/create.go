package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/data-api/dto"
	"net/http"
)

func (controller ServiceController) Create(ctx context.Context, req core.Request) (*core.Response, error) {

	_, webhookStatusCode, webhookError := controller.WebhookService.Pre(ctx, req)
	if webhookError != nil {
		return &core.Response{
			StatusCode: webhookStatusCode,
			Body:       webhookError,
		}, nil
	}

	var body dto.CreateConfig
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

	id, revision, err := controller.DataService.Create(ctx, base, collection, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusCreated,
		Body:       dto.CreateResponse{ID: id, Revision: revision},
	}, nil
}
