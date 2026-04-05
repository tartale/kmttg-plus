package model

import (
	"github.com/tartale/go/pkg/filter"
)

type (
	TivoFilterFn = filter.TypeFilter[*Tivo]
	ShowFilterFn = filter.TypeFilter[*Show]
)
