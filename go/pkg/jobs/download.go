package jobs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/retry"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"go.uber.org/zap"
)

func Download(ctx context.Context, subtask *Subtask) error {

	downloadDir := subtask.OutputDir()
	if filez.IsDir(downloadDir) {
		return subtask.Complete(ctx)
	}
	tmpDir := subtask.Tmpdir()
	err := os.MkdirAll(tmpDir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("error while trying to create directory: %w", err)
	}

	logz.Logger.Debug("started downloading show", zap.String("showID", subtask.ShowID))
	retry.Eventually(func() error {
		subtask.Status.Progress += 10

		return nil
	}, 10*time.Second, 1*time.Second)
	logz.Logger.Debug("finished downloading show", zap.String("showID", subtask.ShowID))

	return subtask.Complete(ctx)
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
