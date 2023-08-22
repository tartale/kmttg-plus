package apicontext

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
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

func (a APIContext) GqlValue(path string, key ContextKey) any {

	if !graphql.HasOperationContext(a) {
		return nil
	}
	fctx := graphql.GetFieldContext(a)
	octx := graphql.GetOperationContext(a)
	if fctx == nil || octx == nil {
		return nil
	}
	if path == "" {
		if val, ok := fctx.Args[string(key)]; ok {
			return val
		}
		return nil
	}

	collectedFields := graphql.CollectFields(octx, fctx.Field.Selections, nil)
	for _, cf := range collectedFields {
		child, err := fctx.Child(a, cf)
		if err != nil {
			return nil
		}
		if child.Path().String() == path {
			if val, ok := child.Args[string(key)]; ok {
				return val
			}
		}
	}

	return nil
}

func ShowOffset(ctx context.Context) int {
	val := ctx.Value(ShowOffsetKey)
	if val != nil {
		return val.(int)
	}
	val = Wrap(ctx).GqlValue("tivos.shows", "offset")
	if val != nil {
		v := val.(*int)
		return *v
	}

	return DefaultOffset
}

func ShowLimit(ctx context.Context) int {
	val := ctx.Value(ShowLimitKey)
	if val != nil {
		return val.(int)
	}
	val = Wrap(ctx).GqlValue("tivos.shows", "limit")
	if val != nil {
		v := val.(*int)
		return *v
	}

	return DefaultLimit
}

func TivoFilterFn(ctx context.Context) model.TivoFilterFn {
	val := ctx.Value(TivoFilterKey)
	if val != nil {
		return val.(model.TivoFilterFn)
	}
	return nil
}

func ShowFilterFn(ctx context.Context) model.ShowFilterFn {
	val := Wrap(ctx).GqlValue("tivos.shows", "filter")
	if val != nil {
		showFilters := val.([]*model.ShowFilter)
		showFilterFn := model.NewShowFilterFn(showFilters)
		return showFilterFn
	}
	return nil
}
