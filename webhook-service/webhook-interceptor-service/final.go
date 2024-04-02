package http_webhook_service

import (
	"context"
	"github.com/devingen/data-api/dto"

	core "github.com/devingen/api-core"
)

func (service WebhookInterceptorService) Final(ctx context.Context, req core.Request, responseBody interface{}) {
	if service.Client == nil {
		return
	}

	service.Client.R().EnableTrace().
		SetBody(dto.WebhookFinalRequest{
			Method:         req.HTTPMethod,
			Path:           req.Path,
			PathParameters: req.PathParameters,
			QueryParams:    req.QueryStringParameters,
			Headers:        req.Headers,
			ResponseBody:   responseBody,
		}).
		Post("/final")
	return
}
