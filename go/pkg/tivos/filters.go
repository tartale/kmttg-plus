package tivos

import (
	"github.com/tartale/go/pkg/filter"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

type FilterFn = func(t *model.Tivo) bool

func NewFilterFn(tf []*model.TivoFilter) FilterFn {
	return func(t *model.Tivo) bool {
		if len(tf) == 0 {
			return true
		}
		expression := filter.GetExpression(tf)
		values := filter.GetValues(tf, t)
		eval := filter.MustEvaluate(expression, values)

		return eval.(bool)
	}
}
