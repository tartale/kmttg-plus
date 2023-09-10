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

	val, err := gqlgen.GetArgValue[[]*model.ShowFilter](ctx, apicontext.ShowFiltersKey)
	if err != nil && errors.Is(err, generics.ErrNotCasted) {
		return nil, fmt.Errorf("%w '%s'; expected type: %s", errorz.ErrInvalidArgument,
			apicontext.ShowFiltersKey.Name, "[ShowFilter]")
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
		eval, err := gval.Evaluate(expression, values)
		if err != nil {
			logz.Logger.Warn("error attempting to filter shows", zap.Error(err))
			return true
		}

		return eval.(bool)
	}
}

func GetImageDimensions(ctx context.Context) (*apicontext.ImageDimensions, error) {
	var result apicontext.ImageDimensions

	height, err := getDimension(ctx, apicontext.ShowImageURLHeightKey)
	if err != nil {
		return nil, err
	}

	width, err := getDimension(ctx, apicontext.ShowImageURLWidthKey)
	if err != nil {
		return nil, err
	}

	if height == nil || width == nil {
		return nil, nil
	}

	if *width <= 0 || *height <= 0 {
		return nil, fmt.Errorf("%w: image height and width must both be positive integers", errorz.ErrInvalidArgument)
	}

	result.Height = *height
	result.Width = *width

	return &result, nil
}

func getDimension(ctx context.Context, key gqlgen.ArgKey) (*int, error) {

	value, err := gqlgen.GetArgValue[int](ctx, key)
	if errors.Is(err, gqlgen.ErrArgumentNotFound) {
		return nil, fmt.Errorf("%w: must provide both image height and width", errorz.ErrBadRequest)
	}
	if err != nil && !errors.Is(err, gqlgen.ErrFieldNotFound) {
		return nil, err
	}

	return value, nil
}
