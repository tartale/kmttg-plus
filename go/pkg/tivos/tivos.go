package tivos

import (
	"context"
	"sort"
	"sync"

	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

var tivos = make(map[string]*model.Tivo)
var tivoClients = make(map[string]*client.TivoClient)
var tivoMutex sync.Mutex

func Add(tivo *model.Tivo) {

	tivoMutex.Lock()
	defer tivoMutex.Unlock()
	tivos[tivo.Name] = tivo

	logz.Logger.Info("updated tivo list", zap.String("name", tivo.Name), zap.String("address", tivo.Address))
}

func GetClient(tivo *model.Tivo) (*client.TivoClient, error) {

	tivoMutex.Lock()
	defer tivoMutex.Unlock()
	if tivoClient, ok := tivoClients[tivo.Name]; ok {
		return tivoClient, nil
	}
	newTivoClient, err := client.NewTivoClient(tivo)
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
