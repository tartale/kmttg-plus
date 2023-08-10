package apicontext

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type ContextKey string

const (
	OffsetKey ContextKey = "offset"
	LimitKey  ContextKey = "limit"
	FilterKey ContextKey = "filter"

	DefaultOffset = 0
	DefaultLimit  = 25
)

type FilterFn func(any) bool

type APIContext struct {
	context.Context
}

func Wrap(ctx context.Context) APIContext {
	return APIContext{
		Context: ctx,
	}
}

func (a APIContext) WithOffset(offset int) APIContext {
	return Wrap(context.WithValue(a, OffsetKey, offset))
}

func (a APIContext) WithLimit(limit int) APIContext {
	return Wrap(context.WithValue(a, LimitKey, limit))
}

func (a APIContext) WithFilter(filter any) APIContext {
	return Wrap(context.WithValue(a, FilterKey, filter))
}

func (a APIContext) GqlValue(path string, key ContextKey) any {

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

func Offset(ctx context.Context) int {
	val := ctx.Value(OffsetKey)
	if val != nil {
		return val.(int)
	}
	return DefaultOffset
}

func Limit(ctx context.Context) int {
	val := ctx.Value(LimitKey)
	if val != nil {
		return val.(int)
	}
	return DefaultLimit
}

func Filter(ctx context.Context) any {
	val := ctx.Value(FilterKey)
	if val != nil {
		return val
	}
	return nil
}
