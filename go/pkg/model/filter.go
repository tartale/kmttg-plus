package model

func Filter[T any](slice []T, f func(T) bool) []T {

	var result []T
	for _, item := range slice {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}

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
