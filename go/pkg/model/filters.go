package model

type TivoFilterFn = func(t *Tivo) bool
type ShowFilterFn = func(s Show) bool
