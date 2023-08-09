package client

import (
	"context"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var clients = xsync.NewMapOf[*TivoClient]()

func Get(tivo *model.Tivo) (*TivoClient, error) {

	if client, ok := clients.Load(tivo.Name); ok {
		return client, nil
	}

	newTivoClient, err := NewTivoClient(tivo)
	if err != nil {
		return nil, err
	}
	err = newTivoClient.Authenticate(context.Background())
	if err != nil {
		return nil, err
	}
	clients.Store(tivo.Name, newTivoClient)

	return newTivoClient, nil
}
