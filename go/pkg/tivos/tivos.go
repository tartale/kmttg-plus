package tivos

import (
	"sort"
	"sync"

	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

var tivos = make(map[string]*model.Tivo)
var tivoMutex sync.Mutex

func Add(tivo *model.Tivo) {

	tivoMutex.Lock()
	defer tivoMutex.Unlock()
	tivos[tivo.Name] = tivo

	logz.Logger.Info("updated tivo list", zap.String("name", tivo.Name))
}

func List() []*model.Tivo {
	list := maps.Values(tivos)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}
