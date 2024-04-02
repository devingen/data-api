package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/data-api/dto"
	"net/http"
)

func (controller ServiceController) Delete(ctx context.Context, req core.Request) (*core.Response, error) {

	_, webhookStatusCode, webhookError := controller.WebhookService.Pre(ctx, req)
	if webhookError != nil {
		return &core.Response{
			StatusCode: webhookStatusCode,
			Body:       webhookError,
		}, nil
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

	_, _, err := controller.DataService.Delete(ctx, base, collection, id)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       dto.DeleteResponse{},
	}, nil
}
