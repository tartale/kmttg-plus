package recordings

import (
	"context"

	"github.com/tartale/kmttg-plus/go/pkg/mindrpc"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func Get(ctx context.Context, tivo *model.Tivo) ([]*model.Show, error) {

	client, err := mindrpc.NewTivoClient(tivo)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// err = client.Request()
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
