package tivo

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

const retryCount = 3

func GetShows(ctx context.Context, tivo *model.Tivo) ([]model.Show, error) {

	var result []model.Show
	var err error

	tivoClient, err := GetClient(tivo)
	if err != nil {
		return nil, err
	}

	for retries := 0; retries < retryCount; retries++ {
		result, err = tivoClient.GetShows(ctx)
		if errors.Is(err, errorz.ErrReconnected) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get shows: %w", err)
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed to get shows; number of retries exceeded: %w", err)
}

func GetEpisodes(ctx context.Context, series *model.Series) ([]*model.Episode, error) {

	tivoClient, err := GetClient(series.Tivo)
	_ = tivoClient
	if err != nil {
		return nil, err
	}

	return nil, nil
}
