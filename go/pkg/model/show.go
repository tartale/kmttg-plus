package model

import (
	"fmt"

	"github.com/tartale/go/pkg/primitive"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
)

func NewShow(recordingItem message.RecordingItem, parent message.RecordingFolderItem) (Show, error) {

	switch recordingItem.CollectionType {
	case message.CollectionTypeSeries:
		var episodeNumber int
		if len(recordingItem.EpisodeNum) > 0 {
			episodeNumber = recordingItem.EpisodeNum[0]
		}
		return &Episode{
			Kind:               ShowKindEpisode,
			RecordingID:        parent.ChildRecordingID,
			Title:              recordingItem.Title,
			RecordedOn:         recordingItem.StartTime.Time,
			Description:        "",
			OriginalAirDate:    recordingItem.OriginalAirDate,
			SeasonNumber:       recordingItem.SeasonNumber,
			EpisodeNumber:      episodeNumber,
			EpisodeTitle:       recordingItem.Subtitle,
			EpisodeDescription: recordingItem.Description,
		}, nil
	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		return &Movie{
			Kind:        ShowKindMovie,
			RecordingID: parent.ChildRecordingID,
			Title:       recordingItem.Title,
			RecordedOn:  recordingItem.StartTime.Time,
			Description: recordingItem.Description,
			MovieYear:   primitive.Ref(recordingItem.MovieYear),
		}, nil
	default:
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingItem.CollectionType)))

	}

}

func MergeEpisodes(shows []Show) []Show {
	combinedShowsMap := make(map[string]Show)

	for _, show := range shows {
		if show.GetKind() == ShowKindEpisode {
			episode := show.(*Episode)
			if existingSeries, exists := combinedShowsMap[episode.Title]; exists {
				series := existingSeries.(*Series)
				series.Episodes = append(series.Episodes, episode)
				if series.RecordedOn.Before(episode.RecordedOn) {
					series.RecordedOn = episode.RecordedOn
				}
			} else {
				combinedShowsMap[episode.Title] = &Series{
					Kind:        ShowKindSeries,
					RecordingID: show.GetRecordingID(),
					Title:       show.GetTitle(),
					RecordedOn:  show.GetRecordedOn(),
					Description: show.GetDescription(),
					Episodes:    []*Episode{episode},
				}
			}
		} else if show.GetKind() == ShowKindMovie {
			combinedShowsMap[show.GetTitle()] = show
		}
	}

	combinedShows := make([]Show, 0, len(combinedShowsMap))
	for _, show := range combinedShowsMap {
		combinedShows = append(combinedShows, show)
	}

	return combinedShows
}

/*

const parseShow = (obj: any): Show => {
  const recording = obj.recording[0]
  const show: Series | Movie = {
    recordingID: recording.recordingID,
    kind: recording.episodic ? ShowKind.Series : ShowKind.Movie,
    title: recording.title,
    recordedOn: new Date(recording.startTime),
    description: recording.description,
    movieYear: recording.movieYear,
    episodes: recording.episodic ? [
      {
        recordingID: recording.recordingID,
        kind: ShowKind.Episode,
        title: recording.title,
        recordedOn: new Date(recording.startTime),
        description: recording.description,
        originalAirDate: new Date(recording.originalAirDate),
        seasonNumber: recording.seasonNumber ? recording.seasonNumber : 0,
        episodeNumber: recording.episodeNum ? recording.episodeNum[0] : 0,
        episodeTitle: recording.subtitle,
        episodeDescription: recording.description,
      }
    ] : []
  }

  return show
}

*/
