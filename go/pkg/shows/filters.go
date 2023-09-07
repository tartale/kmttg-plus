package shows

import (
	"context"
	"errors"
	"fmt"

	"github.com/PaesslerAG/gval"
	"github.com/tartale/go/pkg/filter"
	"github.com/tartale/go/pkg/generics"
	"github.com/tartale/go/pkg/gqlgen"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

func NewFilters(ctx context.Context) ([]*model.ShowFilter, error) {

	val, err := gqlgen.GetArgValueE[[]*model.ShowFilter](ctx, apicontext.ShowFiltersKey)
	if err != nil && errors.Is(err, generics.ErrInvalidType) {
		return nil, fmt.Errorf("%w '%s'; expected type: %s", errorz.ErrInvalidArgument,
			apicontext.ShowFiltersKey.Name, "[ShowFilter]")
	}
	if err != nil && errors.Is(err, gqlgen.ErrNotFound) {
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
		eval, err := gval.Evaluate(expression, values)
		if err != nil {
			logz.Logger.Warn("error attempting to filter shows", zap.Error(err))
			return true
		}

		return eval.(bool)
	}
}
