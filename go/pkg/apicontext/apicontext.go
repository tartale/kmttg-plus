package apicontext

import (
	"context"

	"github.com/tartale/go/pkg/contexts"
	"github.com/tartale/go/pkg/gqlgen"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

type ContextKey string

const (
	TivoFilterKey ContextKey = "tivoFilter"
	ShowFilterKey ContextKey = "showFilter"
	ShowOffsetKey ContextKey = "showOffset"
	ShowLimitKey  ContextKey = "showLimit"

	DefaultOffset = 0
	DefaultLimit  = 25
)

type APIContext struct {
	context.Context
}

func Wrap(ctx context.Context) APIContext {
	return APIContext{
		Context: ctx,
	}
}

func (a APIContext) WithShowOffset(offset int) APIContext {
	return Wrap(context.WithValue(a, ShowOffsetKey, offset))
}

func (a APIContext) WithShowLimit(limit int) APIContext {
	return Wrap(context.WithValue(a, ShowLimitKey, limit))
}

func (a APIContext) WithTivoFilterFn(filter model.TivoFilterFn) APIContext {
	return Wrap(context.WithValue(a, TivoFilterKey, filter))
}

func ShowOffset(ctx context.Context) int {
	if val := contexts.Value[int](ctx, ShowOffsetKey); val != nil {
		return *val
	}
	if val := gqlgen.GetArgValueT[int](ctx, "tivos.shows", "offset"); val != nil {
		return *val
	}

	return DefaultOffset
}

func ShowLimit(ctx context.Context) int {
	if val := contexts.Value[int](ctx, ShowLimitKey); val != nil {
		return *val
	}
	if val := gqlgen.GetArgValueT[int](ctx, "tivos.shows", "limit"); val != nil {
		return *val
	}

	return DefaultLimit
}

func TivoFilterFn(ctx context.Context) model.TivoFilterFn {
	if val := contexts.Value[model.TivoFilterFn](ctx, TivoFilterKey); val != nil {
		return *val
	}

	return nil
}

func ShowFilterFn(ctx context.Context) model.ShowFilterFn {
	if val := gqlgen.GetArgValue(ctx, "tivos.shows", "filter"); val != nil {
		showFilters := val.([]*model.ShowFilter)
		showFilterFn := model.NewShowFilterFn(showFilters)
		return showFilterFn
	}

	return nil
}

func ImageURLWidth(ctx context.Context) int {
	if val := gqlgen.GetArgValueT[int](ctx, "tivos.shows.imageURL", "width"); val != nil {
		return *val
	}

	return 0
}

func ImageURLHeight(ctx context.Context) int {
	if val := gqlgen.GetArgValueT[int](ctx, "tivos.shows.imageURL", "height"); val != nil {
		return *val
	}

	return 0
}
