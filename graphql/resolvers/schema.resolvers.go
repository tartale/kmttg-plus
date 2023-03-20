package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"

	models1 "github.com/tartale/kmttg-plus/graphql/model"
	server1 "github.com/tartale/kmttg-plus/graphql/server"
)

// Tivos is the resolver for the tivos field.
func (r *queryResolver) Tivos(ctx context.Context) ([]*models1.Tivo, error) {
	return r.Resolver.Tivos(ctx)
}

// Query returns server1.QueryResolver implementation.
func (r *Resolver) Query() server1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }