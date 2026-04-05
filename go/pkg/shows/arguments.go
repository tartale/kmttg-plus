package shows

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/go/pkg/errorz"
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
		return nil, fmt.Errorf("image height and width must both be positive integers: %w", errorz.ErrInvalidArgument)
	}

	result.Height = *height
	result.Width = *width

	return &result, nil
}

func getDimension(ctx context.Context, key gqlgen.ArgKey) (*int, error) {
	value, err := gqlgen.GetArgValue[int](ctx, key)
	if errors.Is(err, gqlgen.ErrArgumentNotFound) {
		return nil, fmt.Errorf("must provide both image height and width: %w", errorz.ErrBadRequest)
	}
	if err != nil && !errors.Is(err, gqlgen.ErrFieldNotFound) {
		return nil, err
	}

	return value, nil
}
