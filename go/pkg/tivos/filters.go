package tivos

import (
	"github.com/PaesslerAG/gval"
	"github.com/tartale/go/pkg/filter"
	"github.com/tartale/go/pkg/maps"
	"github.com/tartale/go/pkg/structs"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

type FilterFn = func(t *model.Tivo) bool

func NewFilterFn(tf []*model.TivoFilter) FilterFn {
	return func(t *model.Tivo) bool {
		if len(tf) == 0 {
			return true
		}
		expression := filter.GetExpression(tf)
		structWrapper := structs.New(t)
		structWrapper.TagName = "json"
		values := maps.CastPrimitives(structWrapper.Map())
		eval, err := gval.Evaluate(expression, values)
		if err != nil {
			logz.Logger.Warn("error attempting to filter tivos", zap.Error(err))
			return true
		}

		return eval.(bool)
	}
}
