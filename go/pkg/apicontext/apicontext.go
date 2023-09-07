package apicontext

import (
	"context"

	"github.com/tartale/go/pkg/contexts"
	"github.com/tartale/go/pkg/gqlgen"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

const (
	DefaultOffset = 0
	DefaultLimit  = 25
	DefaultHeight = 0
	DefaultWidth  = 0
)

var (
	TivoFiltersKey        = gqlgen.ArgKey{Path: "tivos", Name: "filters"}
	ShowFiltersKey        = gqlgen.ArgKey{Path: "tivos.shows", Name: "filters"}
	ShowOffsetKey         = gqlgen.ArgKey{Path: "tivos.shows", Name: "offset"}
	ShowLimitKey          = gqlgen.ArgKey{Path: "tivos.shows", Name: "limit"}
	ShowImageURLHeightKey = gqlgen.ArgKey{Path: "tivos.shows.imageURL", Name: "height"}
	ShowImageURLWidthKey  = gqlgen.ArgKey{Path: "tivos.shows.imageURL", Name: "width"}
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

func (a APIContext) WithTivoFilterFn(fn model.TivoFilterFn) APIContext {
	return Wrap(context.WithValue(a, TivoFiltersKey, fn))
}

func (a APIContext) WithShowFilterFn(fn model.ShowFilterFn) APIContext {
	return Wrap(context.WithValue(a, ShowFiltersKey, fn))
}

func ShowOffset(ctx context.Context) int {
	if val := contexts.Value[int](ctx, ShowOffsetKey); val != nil {
		return *val
	}
	if val := gqlgen.GetArgValue[int](ctx, ShowOffsetKey); val != nil {
		return *val
	}

	return DefaultOffset
}

func ShowLimit(ctx context.Context) int {
	if val := contexts.Value[int](ctx, ShowLimitKey); val != nil {
		return *val
	}
	if val := gqlgen.GetArgValue[int](ctx, ShowLimitKey); val != nil {
		return *val
	}

	return DefaultLimit
}

func TivoFilterFn(ctx context.Context) model.TivoFilterFn {
	if val := contexts.Value[model.TivoFilterFn](ctx, TivoFiltersKey); val != nil {
		return *val
	}

	return nil
}

func ShowFilterFn(ctx context.Context) model.ShowFilterFn {
	if val := contexts.Value[model.ShowFilterFn](ctx, ShowFiltersKey); val != nil {
		return *val
	}

	return nil
}

func ShowImageURLHeight(ctx context.Context) int {
	if val := gqlgen.GetArgValue[int](ctx, ShowImageURLHeightKey); val != nil {
		return *val
	}

	return DefaultHeight
}

func ShowImageURLWidth(ctx context.Context) int {
	if val := gqlgen.GetArgValue[int](ctx, ShowImageURLWidthKey); val != nil {
		return *val
	}

	return DefaultWidth
}
