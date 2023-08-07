package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"

	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/server"
	"github.com/tartale/kmttg-plus/go/pkg/tivo"
)

// Tivos is the resolver for the tivos field.
func (r *queryResolver) Tivos(ctx context.Context, filter *model.TivoFilter) ([]*model.Tivo, error) {
	tivoFilter := model.NewTivoFilter(filter)
	ctx = apicontext.New(ctx).WithFilter(tivoFilter)
	return tivo.List(ctx), nil
}

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
