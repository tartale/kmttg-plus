package apicontext

import (
	"context"
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

func New(parent context.Context) APIContext {
	return APIContext{parent}
}

func (a APIContext) WithOffset(offset int) APIContext {
	return APIContext{context.WithValue(a.Context, OffsetKey, offset)}
}

func (a APIContext) WithLimit(limit int) APIContext {
	return APIContext{context.WithValue(a.Context, LimitKey, limit)}
}

func (a APIContext) WithFilter(filter any) APIContext {
	return APIContext{context.WithValue(a.Context, FilterKey, filter)}
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
