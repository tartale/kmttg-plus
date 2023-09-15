package jobs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/httpx"
	"github.com/tartale/go/pkg/primitives"
	"github.com/tartale/go/pkg/stringz"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
)

func Download(ctx context.Context, subtask *Subtask) error {

	_, downloadPath := getDownloadPaths(subtask)
	if filez.Exists(downloadPath) {
		subtask.Status.Progress = 100
		return nil
	}

	return download(ctx, subtask)
}

func getDownloadURL(show model.Show) (*url.URL, error) {

	showDetails := shows.GetDetails(show)
	showTitle := url.PathEscape(show.GetTitle())
	showIDNumber := shows.ParseIDNumber(showDetails.ObjectID)
	downloadBaseURL := fmt.Sprintf("http://%s/download/%s.Tivo",
		showDetails.Tivo.Address, showTitle)
	downloadURL, err := url.Parse(downloadBaseURL)
	if err != nil {
		return nil, err
	}
	downloadURL.RawQuery = url.Values{
		"Container": {"NowPlaying"},
		"Format":    {"video/x-tivo-mpeg-ts"},
		"id":        {showIDNumber},
	}.Encode()

	return downloadURL, nil
}

func getDownloadPaths(subtask *Subtask) (tmpPath, outputPath string) {

	tmpPath = path.Join(subtask.tmpdir, stringz.ToAlphaNumeric(subtask.show.GetTitle())+".ts.tmp")
	outputPath = path.Join(subtask.outputdir, stringz.ToAlphaNumeric(subtask.show.GetTitle())+".ts")

	return
}

func download(ctx context.Context, subtask *Subtask) error {

	downloadURL, err := getDownloadURL(subtask.show)
	if err != nil {
		return fmt.Errorf("%w: unable get download URL", err)
	}
	client, err := client.NewHttpClient(shows.GetDetails(subtask.show).Tivo)
	if err != nil {
		return err
	}
	req, err := client.NewRequestWithContext(ctx, http.MethodGet, downloadURL.String(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return httpx.GetResponseError(resp)
	}

	tmpPath, downloadPath := getDownloadPaths(subtask)
	tmpFile := filez.MustOpenFile(tmpPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer tmpFile.Close()

	estimatedLength, err := primitives.ParseTo[int](resp.Header.Get("TiVo-Estimated-Length"))
	if err != nil {
		return err
	}
	progressWriter := NewProgressWriter(subtask, int64(estimatedLength))
	_, err = io.Copy(io.MultiWriter(tmpFile, progressWriter), resp.Body)
	if err != nil {
		return err
	}
	filez.MustRename(tmpPath, downloadPath)

	return nil
}
