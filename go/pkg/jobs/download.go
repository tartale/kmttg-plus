package jobs

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
)

func Download(ctx context.Context, subtask *Subtask) error {

	downloadDir := subtask.OutputDir()
	if filez.IsDir(downloadDir) {
		subtask.Status.Progress = 100
		return nil
	}
	tmpDir := subtask.Tmpdir()
	err := os.MkdirAll(tmpDir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("%w: unable to create directory '%s'", err, tmpDir)
	}
	dlURL, err := getDownloadURL(subtask.show)
	if err != nil {
		return fmt.Errorf("%w: unable get download URL '%s'", err, downloadDir)
	}
	fmt.Println(dlURL)
	subtask.Status.Progress = 100

	err = os.MkdirAll(downloadDir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("%w: unable to create directory '%s'", err, downloadDir)
	}

	return nil
}

func getDownloadURL(show model.Show) (*url.URL, error) {

	showDetails := shows.GetDetails(show)
	showTitle := url.PathEscape(show.GetTitle())
	showIDNumber := shows.ParseIDNumber(showDetails.ObjectaID)
	downloadBaseURL := fmt.Sprintf("http://tivo:%s@%s/download/%s.Tivo",
		config.Values.MediaAccessKey, showDetails.Tivo.Address, showTitle)
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

/*
   const recordingMeta = await this.sendRequest({
       type: 'idSearch',
       objectId: recordingId,
       namespace: 'mfs',
   });
   console.log(recordingMeta);

   const downloadId = recordingMeta.objectId[0].replace('mfs:rc.', '');
   const dUrl = new URL('http://localhost/download/download.TiVo?Container=%2FNowPlaying');
   dUrl.password = this.mak
   dUrl.username = 'tivo';

   dUrl.host = this.ip;
   dUrl.searchParams.append('id', downloadId);
   useTs && dUrl.searchParams.append('Format','video/x-tivo-mpeg-ts');
   return dUrl.toString();

    `http://${this.ip}/download/download.TiVo?Container=%2FNowPlaying&id=` + encodeURIComponent(downloadId) + (useTs ? '&Format=video/x-tivo-mpeg-ts' : '');

*/
