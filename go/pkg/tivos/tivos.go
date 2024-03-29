package tivos

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorx"
	liberrorz "github.com/tartale/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
	"go.uber.org/zap"
)

var (
	tivoMap = xsync.NewMapOf[*model.Tivo]()
)

func RunBackgroundLoader(ctx context.Context) {
	loadTicker := time.NewTicker(10 * time.Second)

	for range loadTicker.C {
		err := LoadAll()
		if err != nil {
			logz.Logger.Warn("error loading shows", zap.Error(err))
			loadTicker.Reset(30 * time.Second)
		} else {
			loadTicker.Reset(5 * time.Minute)
		}
	}
}

func LoadAll() error {

	var errs errorx.Errors
	tivoMap.Range(func(key string, val *model.Tivo) bool {
		errs = append(errs, Load(val))
		return true
	})

	return errs.Combine("errors when loading shows", "\n")
}

func Load(tivo *model.Tivo) error {

	logz.Logger.Debug("loading all shows", zap.String("tivoName", tivo.Name))
	tivoClient, err := client.NewRpcClient(tivo)
	if err != nil {
		return err
	}

	shows, err := LoadShows(context.Background(), tivoClient)
	if err != nil {
		return err
	}

	newTivo := *tivo
	newTivo.Shows = shows
	tivoMap.Store(tivo.Name, &newTivo)
	logz.Logger.Debug("Successfully loaded all shows", zap.String("tivoName", tivo.Name))

	return nil
}

func LoadShows(ctx context.Context, tivoClient *client.TivoClient) ([]model.Show, error) {

	const (
		retryCount = 3
	)

	var (
		shows   []model.Show
		success bool
		err     error
	)

	for retries := 0; retries < retryCount; retries++ {
		shows, err = tivoClient.GetShows(ctx)
		if errors.Is(err, errorz.ErrReconnected) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get shows: %w", err)
		}

		success = true
		break
	}
	if !success {
		return nil, fmt.Errorf("failed to get shows; number of retries exceeded: %w", err)
	}

	return shows, nil
}

func List(ctx context.Context) []*model.Tivo {

	var list []*model.Tivo
	tivoFilterFn := apicontext.TivoFilterFn(ctx)
	showFilterFn := apicontext.ShowFilterFn(ctx)
	imageDimensions := apicontext.ShowImageDimensions(ctx)

	tivoMap.Range(func(key string, val *model.Tivo) bool {
		if tivoFilterFn != nil && !tivoFilterFn(val) {
			return true
		}
		tivo := *val
		list = append(list, &tivo)

		tivo.Shows = []model.Show{}
		offsetCountdown := apicontext.ShowOffset(ctx)
		limitCountdown := apicontext.ShowLimit(ctx)
		for _, show := range val.Shows {
			if limitCountdown == 0 {
				break
			}
			if offsetCountdown > 0 {
				offsetCountdown--
				continue
			}
			if showFilterFn != nil && !showFilterFn(show) {
				continue
			}
			show = shows.WithImageURL(show, imageDimensions)
			tivo.Shows = append(tivo.Shows, shows.AsAPIType(show))
			limitCountdown--
		}

		return true
	})

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

func GetShowForID(recordingID string) (model.Show, error) {

	var result model.Show
	tivoMap.Range(func(key string, val *model.Tivo) bool {

		for _, show := range val.Shows {
			details := shows.GetDetails(show)
			if details.Recording.RecordingID == recordingID {
				clone := shows.New(val, details.ObjectID, &details.Recording, &details.Collection)
				result = clone
				return false
			}
		}

		return true
	})

	if result == nil {
		return nil, fmt.Errorf("%w: show ID '%s'", liberrorz.ErrNotFound, recordingID)
	}

	return result, nil
}
