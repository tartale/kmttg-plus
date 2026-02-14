package tivos

import (
	"github.com/PaesslerAG/gval"
	"github.com/tartale/go/pkg/filter"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

type FilterFn = func(t *model.Tivo) bool

func NewFilterFn(f []*model.TivoFilter) FilterFn {

	return func(t *model.Tivo) bool {

		if len(f) == 0 {
			return true
		}
		expression := filter.GetExpression(f)
		values := filter.GetValues(f, t)
		eval, err := gval.Evaluate(expression, values)
		if err != nil {
			logz.Logger.Warn("error attempting to filter tivos", zap.Error(err))
			return true
		}

		return eval.(bool)
	}
}
