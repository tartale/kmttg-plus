package client

import (
	"context"
	"encoding/json"
	"github.com/tartale/kmttg-plus/graphql/model"
	"io"
	"net/http"
)

const kmttgURL = "http://10.0.1.3:8181"

func Tivos(ctx context.Context) ([]*model.Tivo, error) {
	result, err := get[[]*model.Tivo](ctx, kmttgURL+"/getTivos")
	if err != nil {
		return nil, err
	}

	return *result, nil
}

func get[T any](ctx context.Context, url string) (*T, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal(body, &result)

	return &result, nil
}
