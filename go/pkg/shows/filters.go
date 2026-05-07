package shows

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/filter"
	"github.com/tartale/go/pkg/generics"
	"github.com/tartale/go/pkg/gqlgen"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func NewFilters(ctx context.Context) ([]*model.ShowFilter, error) {
	val, err := gqlgen.GetArgValue[[]*model.ShowFilter](ctx, apicontext.ShowFiltersKey)
	if err != nil && errors.Is(err, generics.ErrNotCasted) {
		return nil, fmt.Errorf("'%s'; expected type: %s: %w", apicontext.ShowFiltersKey.Name, "[ShowFilter]", errorz.ErrInvalidArgument)
	}
	if err != nil && errors.Is(err, errorz.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return *val, nil
}

func NewFilterFn(sf []*model.ShowFilter) model.ShowFilterFn {
	return func(s model.Show) bool {
		if len(sf) == 0 {
			return true
		}
		expression := filter.GetExpression(sf)
		values := filter.GetValues(sf, s)
		eval := filter.MustEvaluate(expression, values)

		return eval.(bool)
	}
}
