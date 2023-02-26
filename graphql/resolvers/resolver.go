package resolvers

import (
	"context"
	"github.com/tartale/kmttg-plus/graphql/client"
	"github.com/tartale/kmttg-plus/graphql/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func (r Resolver) Tivos(ctx context.Context) ([]*model.Tivo, error) {
	return client.Tivos(ctx)
}
