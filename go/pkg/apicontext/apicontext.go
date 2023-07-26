package apicontext

import (
	"context"
)

type ContextKey string

const (
	OffsetKey ContextKey = "offset"
	LimitKey  ContextKey = "limit"

	DefaultOffset = 0
	DefaultLimit  = 25
)

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
