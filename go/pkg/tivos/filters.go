package tivos

import (
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func NewFilterFn(tf model.TivoFilters) model.TivoFilterFn {
	return func(t *model.Tivo) bool {
		return tf.ShouldInclude(t)
	}
}
