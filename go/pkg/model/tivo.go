package model

import "fmt"

type TivoFilterFn = func(t *Tivo) bool

func NewTivoFilter(f ...*TivoFilter) TivoFilterFn {

	return func(t *Tivo) bool {
		if len(f) == 0 {
			return true
		}
		if len(f) != 1 {
			panic(fmt.Errorf("only one filter is allowed; received: %d", len(f)))
		}
		if f[0] == nil {
			return true
		}
		if f[0].Name != nil && t.Name == *f[0].Name {
			return true
		}

		return false
	}
}
