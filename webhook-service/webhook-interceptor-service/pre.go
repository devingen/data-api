package http_webhook_service

import (
	"context"
	"encoding/json"
	"github.com/devingen/data-api/dto"
	"net/http"
	"net/url"

	core "github.com/devingen/api-core"
)

func (service WebhookInterceptorService) Pre(ctx context.Context, req core.Request) (*dto.WebhookPreResponse, int, interface{}) {
	if service.Client == nil {
		return nil, 0, nil
	}

	body := make(map[string]interface{})

	if req.Body != "" {
		err := json.Unmarshal([]byte(req.Body), &body)
		if err != nil {
			return nil, 0, err
		}
	}

	resp, err := service.Client.R().EnableTrace().
		SetBody(dto.WebhookPreRequest{
			Method:      req.HTTPMethod,
			Path:        req.Path,
			QueryParams: req.QueryStringParameters,
			Headers:     req.Headers,
			Body:        body,
		}).
		SetResult(&dto.WebhookPreResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/pre")
	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "webhook-api-is-unreachable")
		}
		return nil, resp.StatusCode(), err
	}
	if resp.StatusCode() > 399 {
		return nil, resp.StatusCode(), resp.Error()
	}

	webhookResponse := resp.Result().(*dto.WebhookPreResponse)
	return webhookResponse, resp.StatusCode(), nil
}
