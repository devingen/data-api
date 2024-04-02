package controller

import (
	"context"
	core "github.com/devingen/api-core"
)

type IServiceController interface {
	Query(ctx context.Context, req core.Request) (*core.Response, error)
	Create(ctx context.Context, req core.Request) (*core.Response, error)
	Update(ctx context.Context, req core.Request) (*core.Response, error)
	Delete(ctx context.Context, req core.Request) (*core.Response, error)
}
