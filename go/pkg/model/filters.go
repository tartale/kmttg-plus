package model

import "github.com/tartale/go/pkg/filter"

type (
	TivoFilterFn = func(t *Tivo) bool
	ShowFilterFn = func(s Show) bool
)

type TivoFilters []*TivoFilter

func (f TivoFilters) ShouldInclude(v *Tivo) bool {
	if len(f) == 0 {
		return true
	}
	return filter.ShouldInclude(f, v)
}

type ShowFilters []*ShowFilter

func (f ShowFilters) ShouldInclude(v Show) bool {
	if len(f) == 0 {
		return true
	}
	return filter.ShouldInclude(f, v)
}
