package apicontext

import (
	"context"
	"fmt"

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

func (a APIContext) GqlArgValue(path string, argName string) any {

	fctx := graphql.GetFieldContext(a)
	field := a.findGqlField(path, fctx)
	if field == nil {
		return nil
	}
	arg := field.Arguments.ForName(argName)
	if arg == nil {
		return nil
	}
	val, err := arg.Value.Value(nil)
	if err != nil {
		return nil
	}

	return val
}

func (a APIContext) findGqlField(path string, fctx *graphql.FieldContext) *graphql.CollectedField {

	if !graphql.HasOperationContext(a) {
		return nil
	}
	octx := graphql.GetOperationContext(a)
	if fctx == nil || octx == nil {
		return nil
	}
	fieldPath := fctx.Path()
	fieldPathStr := fieldPath.String()
	if fieldPathStr == path {
		return &fctx.Field
	}

	collectedFields := graphql.CollectFields(octx, fctx.Field.Selections, nil)
	for _, cf := range collectedFields {
		fullPath := fmt.Sprintf("%s.%s", fieldPathStr, cf.Field.Name)
		if path == fullPath {
			return &cf
		}
		child, err := fctx.Child(a, cf)
		if err != nil {
			continue
		}
		childField := a.findGqlField(path, child)
		if childField != nil {
			return childField
		}
	}

	return nil
}

func ShowOffset(ctx context.Context) int {
	val := ctx.Value(ShowOffsetKey)
	if val != nil {
		return val.(int)
	}
	val = Wrap(ctx).GqlArgValue("tivos.shows", "offset")
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
	val = Wrap(ctx).GqlArgValue("tivos.shows", "limit")
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
	val := Wrap(ctx).GqlArgValue("tivos.shows", "filter")
	if val != nil {
		showFilters := val.([]*model.ShowFilter)
		showFilterFn := model.NewShowFilterFn(showFilters)
		return showFilterFn
	}
	return nil
}

func ImageURLWidth(ctx context.Context) int {
	val := Wrap(ctx).GqlArgValue("tivos.shows.imageURL", "width")
	if val == nil {
		return 0
	}

	return int(val.(int64))
}

func ImageURLHeight(ctx context.Context) int {
	val := Wrap(ctx).GqlArgValue("tivos.shows.imageURL", "height")
	if val == nil {
		return 0
	}

	return int(val.(int64))
}
