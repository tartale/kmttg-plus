package tivos

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorx"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/filter"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

const (
	retryCount = 3
)

var (
	tivoMap = xsync.NewMapOf[*model.Tivo]()
)

type TivoFilterFn = func(t *model.Tivo) bool

func NewTivoFilter(f *model.TivoFilter) TivoFilterFn {

	return func(t *model.Tivo) bool {

		expression := filter.GetExpression(f)
		values := filter.GetValues(f, t)
		eval, err := gval.Evaluate(expression, values)
		if err != nil {
			logz.Logger.Warn("error attempting to filter tivos", zap.Error(err))
			return true
		}

		return eval.(bool)
	}
}

func RunBackgroundLoader() {
	loadTicker := time.NewTicker(5 * time.Minute)

	for range loadTicker.C {
		err := LoadAll()
		if err != nil {
			logz.Logger.Warn("error loading shows", zap.Error(err))
			loadTicker.Reset(1 * time.Minute)
		} else {
			loadTicker.Reset(5 * time.Minute)
		}
	}
}

func LoadAll() error {

	var errs errorx.Errors
	tivoz := List(context.Background())
	for _, tvo := range tivoz {
		errs = append(errs, Load(tvo))
	}

	return errs.Combine("errors when loading shows", "\n")
}

func Load(tvo *model.Tivo) error {

	logz.Logger.Debug("loading all shows", zap.String("tivoName", tvo.Name))
	tivoClient, err := client.Get(tvo)
	if err != nil {
		return err
	}

	shows, err := getShows(context.Background(), tivoClient)
	if err != nil {
		return err
	}

	tvo.Shows = shows
	tivoMap.Store(tvo.Name, tvo)
	logz.Logger.Debug("Successfully loaded all shows", zap.String("tivoName", tvo.Name))

	return nil
}

func List(ctx context.Context) []*model.Tivo {

	var list []*model.Tivo
	filterFn := apicontext.Filter(ctx)

	tivoMap.Range(func(key string, val *model.Tivo) bool {
		if filterFn == nil || filterFn.(TivoFilterFn)(val) {
			list = append(list, val)
			return true
		}
		return true
	})

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

func getShows(ctx context.Context, tivoClient *client.TivoClient) ([]model.Show, error) {

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
