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

func NewFilters(ctx context.Context) (model.ShowFilters, error) {
	val, err := gqlgen.GetArgValue[model.ShowFilters](ctx, apicontext.ShowFiltersKey)
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

func NewFilterFn(sf model.ShowFilters) model.ShowFilterFn {
	return func(s model.Show) bool {
		return sf.ShouldInclude(s)
	}
}
