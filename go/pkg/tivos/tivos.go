package tivos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorx"
	liberrorz "github.com/tartale/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
	"go.uber.org/zap"
)

var (
	tivoMap           = xsync.NewMapOf[*model.Tivo]()
	loadFromCacheOnce sync.Once
)

func RunBackgroundLoader(ctx context.Context) {
	loadTicker := time.NewTicker(10 * time.Second)

	for range loadTicker.C {
		err := LoadAll()
		if err != nil {
			logz.Logger.Warn("Error loading shows", zap.Error(err))
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
	loadedFromCache := false

	// Check if there's a cache to initialize from the first time
	loadFromCacheOnce.Do(func() {
		loadedFromCache = loadFromCache(tivo)
	})
	if loadedFromCache {
		return nil
	}

	logz.Logger.Debug("Loading shows via RPC", zap.String("tivoName", tivo.Name))
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
	logz.Logger.Debug("Successfully loaded all shows via RPC", zap.String("tivoName", tivo.Name))
	storeToCache(&newTivo)

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
		return nil, fmt.Errorf("show ID '%s': %w", recordingID, liberrorz.ErrNotFound)
	}

	return result, nil
}

func loadFromCache(tivo *model.Tivo) bool {
	tivoCacheFile := path.Join(config.Values.CacheDir, tivo.Name+".json")
	if _, err := os.Stat(tivoCacheFile); errors.Is(err, os.ErrNotExist) {
		logz.Logger.Debug("No cache found", zap.String("tivoName", tivo.Name))
		return false
	}
	logz.Logger.Debug("Loading shows from cache", zap.String("tivoName", tivo.Name))
	data, err := os.ReadFile(tivoCacheFile)
	if err != nil {
		logz.Logger.Debug("Unable to load cache file", zap.String("tivoName", tivo.Name), zap.Error(err))
		return false
	}
	var aux struct {
		Name    string            `json:"name"`
		Address string            `json:"address"`
		Tsn     string            `json:"tsn"`
		Shows   []json.RawMessage `json:"shows,omitempty"`
	}
	err = json.Unmarshal(data, &aux)
	if err != nil {
		logz.Logger.Debug("Unable to load cache file", zap.String("tivoName", tivo.Name), zap.Error(err))
		return false
	}
	newTivo := model.Tivo{
		Name:    aux.Name,
		Address: aux.Address,
		Tsn:     aux.Tsn,
		Shows:   make([]model.Show, 0, len(aux.Shows)),
	}
	// Unmarshal each show using the shows package helper to handle wrapper types
	for _, raw := range aux.Shows {
		show, err := shows.UnmarshalShowFromJSON(raw, &newTivo)
		if err != nil {
			logz.Logger.Debug("Unable to unmarshal show from cache", zap.String("tivoName", tivo.Name), zap.Error(err))
			return false
		}
		newTivo.Shows = append(newTivo.Shows, show)
	}
	logz.Logger.Debug("Successfully loaded all shows from cache", zap.String("tivoName", tivo.Name))
	tivoMap.Store(tivo.Name, &newTivo)
	return true
}

func storeToCache(tivo *model.Tivo) {
	tivoCacheFile := path.Join(config.Values.CacheDir, tivo.Name+".json")
	data, err := json.MarshalIndent(tivo, "", "  ")
	if err != nil {
		logz.Logger.Debug("Unable to marshal Tivo to JSON; skipping cache write", zap.String("tivoName", tivo.Name), zap.Error(err))
		return
	}
	err = os.WriteFile(tivoCacheFile, data, 0o664)
	if err != nil {
		logz.Logger.Debug("Unable to write cache file", zap.String("tivoName", tivo.Name), zap.Error(err))
		return
	}
	logz.Logger.Debug("Successfully stored all shows to cache", zap.String("tivoName", tivo.Name))
}
