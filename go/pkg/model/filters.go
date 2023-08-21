package model

import (
	"github.com/PaesslerAG/gval"
	"github.com/tartale/kmttg-plus/go/pkg/filter"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"go.uber.org/zap"
)

type TivoFilterFn = func(t *Tivo) bool

func NewTivoFilterFn(f []*TivoFilter) TivoFilterFn {

	return func(t *Tivo) bool {

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

type ShowFilterFn = func(s Show) bool

func NewShowFilterFn(sf []*ShowFilter) ShowFilterFn {

	return func(s Show) bool {

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
