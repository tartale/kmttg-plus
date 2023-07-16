package tivos

import (
	"context"
	"sort"
	"sync"

	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/mindrpc"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

var tivos = make(map[string]*model.Tivo)
var tivoClients = make(map[string]*mindrpc.TivoClient)
var tivoMutex sync.Mutex

func Add(tivo *model.Tivo) {

	tivoMutex.Lock()
	defer tivoMutex.Unlock()
	tivos[tivo.Name] = tivo

	logz.Logger.Info("updated tivo list", zap.String("name", tivo.Name))
}

func GetClient(tivo *model.Tivo) (*mindrpc.TivoClient, error) {

	if tivoClient, ok := tivoClients[tivo.Name]; ok {
		return tivoClient, nil
	}
	tivoMutex.Lock()
	defer tivoMutex.Unlock()
	newTivoClient, err := mindrpc.NewTivoClient(tivo)
	if err != nil {
		return nil, err
	}
	err = newTivoClient.Authenticate(context.Background())
	if err != nil {
		return nil, err
	}

	tivoClients[tivo.Name] = newTivoClient

	return newTivoClient, nil
}

func List() []*model.Tivo {
	list := maps.Values(tivos)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}
