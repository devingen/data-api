package webhook_service

import (
	"context"
	"github.com/devingen/data-api/dto"

	core "github.com/devingen/api-core"
)

// IDataWebhookService defines the functionality of the interceptors
type IDataWebhookService interface {
	Pre(ctx context.Context, req core.Request) (*dto.WebhookPreResponse, int, interface{})
	Final(ctx context.Context, req core.Request, responseBody interface{})
}
