package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/icholy/digest"
	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/httpx"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"golang.org/x/exp/slices"
)

var clients = xsync.NewMapOf[*TivoClient]()

func NewRpcClient(tivo *model.Tivo) (*TivoClient, error) {

	if client, ok := clients.Load(tivo.Name); ok {
		return client, nil
	}

	newRpcClient, err := NewTivoClient(tivo)
	if err != nil {
		return nil, err
	}
	err = newRpcClient.Authenticate(context.Background())
	if err != nil {
		return nil, err
	}
	clients.Store(tivo.Name, newRpcClient)

	return newRpcClient, nil
}

type HttpClient struct {
	http.Client
	sessionCookie *http.Cookie
}

func NewHttpClient(tivo *model.Tivo) (*HttpClient, error) {

	newHttpClient := HttpClient{
		Client: http.Client{
			Transport: &digest.Transport{
				Username: "tivo",
				Password: config.Values.MediaAccessKey,
			},
		},
	}

	tivoHttpAddress := fmt.Sprintf("http://%s", tivo.Address)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodOptions, tivoHttpAddress, nil)
	if err != nil {
		return nil, err
	}
	resp, err := newHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, httpx.GetResponseError(resp)
	}
	sidCookieIndex := slices.IndexFunc(resp.Cookies(), func(c *http.Cookie) bool {
		return c.Name == "sid"
	})
	if sidCookieIndex < 0 {
		return nil, fmt.Errorf("%w: session ID cookie not found", errorz.ErrResponse)
	}
	newHttpClient.sessionCookie = resp.Cookies()[sidCookieIndex]

	return &newHttpClient, nil
}

func (c *HttpClient) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.AddCookie(c.sessionCookie)

	return req, nil
}
