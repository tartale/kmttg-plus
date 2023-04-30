package beacon

import (
	"context"
	"errors"
	"strings"

	"github.com/grandcat/zeroconf"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/tivos"
)

func Listen(ctx context.Context) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)
	if err != nil {
		return err
	}

	for {
		select {
		case entry := <-entries:
			tivo, err := newTivoFromServiceEntry(entry)
			if err != nil {
				continue
			}
			tivos.Add(tivo)
		case <-ctx.Done():
			return nil
		}
	}
}

func newTivoFromServiceEntry(entry *zeroconf.ServiceEntry) (*model.Tivo, error) {

	properties := make(map[string]string)
	for _, property := range entry.Text {
		kv := strings.SplitN(property, "=", 2)
		if len(kv) != 2 {
			continue
		}
		properties[kv[0]] = kv[1]
	}

	var (
		ok            bool
		tsn, platform string
	)

	if tsn, ok = properties["TSN"]; !ok {
		return nil, errors.New("device does not have a TSN")
	}
	if strings.HasPrefix(tsn, "A94") {
		return nil, errors.New("device is not a Tivo DVR")
	}
	if platform, ok = properties["platform"]; !ok {
		return nil, errors.New("device does not have a platform")
	}
	if strings.Contains(platform, "Silver") {
		return nil, errors.New("device is not a Tivo DVR")
	}
	if len(entry.AddrIPv4) == 0 {
		return nil, errors.New("device does not have an IP address")
	}
	name := strings.ReplaceAll(entry.Instance, "\\ ", " ")

	return &model.Tivo{
		Name:    name,
		Address: entry.AddrIPv4[0].String(),
		Tsn:     tsn,
	}, nil
}
