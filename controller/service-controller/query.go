package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/model"
	"github.com/devingen/data-api/dto"
	"net/http"
)

func (controller ServiceController) Query(ctx context.Context, req core.Request) (*core.Response, error) {

	_, webhookStatusCode, webhookError := controller.WebhookService.Pre(ctx, req)
	if webhookError != nil {
		return &core.Response{
			StatusCode: webhookStatusCode,
			Body:       webhookError,
		}, nil
	}

	var body model.QueryConfig
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

	results, meta, err := controller.DataService.Query(ctx, base, collection, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       dto.QueryResponse{Results: results, Meta: meta},
	}, nil
}
