package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/jobs"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/server"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
	"github.com/tartale/kmttg-plus/go/pkg/tivos"
)

// StartJob is the resolver for the startJob field.
func (r *mutationResolver) StartJob(ctx context.Context, job model.Job) (*model.JobStatus, error) {

	newJob := jobs.NewJob(&job)
	jobStatus, err := jobs.StartJob(newJob)
	return jobStatus, err
}

// Tivos is the resolver for the tivos field.
func (r *queryResolver) Tivos(ctx context.Context, filters []*model.TivoFilter) ([]*model.Tivo, error) {

	tivoFilterFn := tivos.NewFilterFn(filters)
	ctx = apicontext.Wrap(ctx).WithTivoFilterFn(tivoFilterFn)

	showFilters, err := shows.NewFilters(ctx)
	if err != nil {
		return nil, err
	}
	showFilterFn := shows.NewFilterFn(showFilters)
	ctx = apicontext.Wrap(ctx).WithShowFilterFn(showFilterFn)

	showImageDimensions, err := shows.GetImageDimensions(ctx)
	if err != nil {
		return nil, err
	}
	ctx = apicontext.Wrap(ctx).WithShowImageDimensions(showImageDimensions)

	return tivos.List(ctx), nil
}

// Jobs is the resolver for the jobs field.
func (r *queryResolver) Jobs(ctx context.Context, filters []*model.JobFilter) ([]*model.JobStatus, error) {

	return jobs.List(ctx, filters)
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
