package shows

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/tivos"
)

const retryCount = 3

func GetRecordingList(ctx context.Context, tivo *model.Tivo) ([]model.Show, error) {

	var result []model.Show
	var err error

	tivoClient, err := tivos.GetClient(tivo)
	if err != nil {
		return nil, err
	}

	for retries := 0; retries < retryCount; retries++ {
		result, err = tivoClient.GetRecordingList(ctx)
		if errors.Is(err, errorz.ErrReconnected) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get recordings: %w", err)
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed to get recordings; number of retries exceeded: %w", err)
}

func GetEpisodes(ctx context.Context, series *model.Series) ([]*model.Episode, error) {

	tivoClient, err := tivos.GetClient(series.Tivo)
	_ = tivoClient
	if err != nil {
		return nil, err
	}

	return nil, nil
}
