package recordings

import (
	"context"

	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/tivos"
)

func Get(ctx context.Context, tivo *model.Tivo) ([]*model.Show, error) {

	tivoClient, err := tivos.GetClient(tivo)
	if err != nil {
		return nil, err
	}

	return tivoClient.GetRecordings(ctx)
}
