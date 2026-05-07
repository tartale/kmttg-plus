package shows

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/gqlgen"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
)

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
