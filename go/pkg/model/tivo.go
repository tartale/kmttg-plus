package model

type TivoFilterFn = func(t *Tivo) bool

func NewTivoFilter(f *TivoFilter) TivoFilterFn {

	return func(t *Tivo) bool {
		if f == nil {
			return true
		}
		if f.Name != nil && t.Name == *f.Name {
			return true
		}

		return false
	}
}
