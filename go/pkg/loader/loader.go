package loader

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorx"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

var Tivos = xsync.NewMapOf[*model.Tivo]()

const retryCount = 3

func LoadAll() error {

	var errs errorx.Errors
	tivoz := List(context.Background())
	for _, tvo := range tivoz {
		errs = append(errs, LoadTivo(tvo))
	}

	return errs.Combine("errors when loading shows", "\n")
}

func LoadTivo(tvo *model.Tivo) error {

	logz.Logger.Debug("loading shows", zap.String("tivoName", tvo.Name))
	tivoClient, err := client.Get(tvo)
	if err != nil {
		return err
	}
	err = PopulateShows(context.Background(), tivoClient, tvo)
	if err != nil {
		return err
	}
	Tivos.Store(tvo.Name, tvo)

	return nil
}

func List(ctx context.Context) []*model.Tivo {

	var list []*model.Tivo
	filterFn := apicontext.Filter(ctx)

	Tivos.Range(func(key string, val *model.Tivo) bool {
		if filterFn == nil || filterFn.(model.TivoFilterFn)(val) {
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

func PopulateShows(ctx context.Context, tivoClient *client.TivoClient, tivo *model.Tivo) error {

	var shows []model.Show
	var success bool
	var err error

	for retries := 0; retries < retryCount; retries++ {
		shows, err = tivoClient.GetShows(ctx)
		if errors.Is(err, errorz.ErrReconnected) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to get shows: %w", err)
		}
		success = true
		break
	}
	if !success {
		return fmt.Errorf("failed to get shows; number of retries exceeded: %w", err)
	}

	tivo.Shows = shows

	return nil
}

func Run() {
	loadTicker := time.NewTicker(30 * time.Second)

	for range loadTicker.C {
		err := LoadAll()
		if err != nil {
			logz.Logger.Warn("error loading shows", zap.Error(err))
		}
	}
}
