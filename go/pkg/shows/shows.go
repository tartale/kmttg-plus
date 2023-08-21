package shows

import (
	"fmt"
	"sort"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func New(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (model.Show, error) {

	switch recordingDetails.CollectionType {
	case message.CollectionTypeSeries:
		return NewEpisode(recordingDetails, collectionDetails)

	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		return NewMovie(recordingDetails)

	default:
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

}

func NewMovie(recordingDetails *message.RecordingItem) (*model.Movie, error) {
	if recordingDetails.CollectionType != message.CollectionTypeMovie && recordingDetails.CollectionType != message.CollectionTypeSpecial {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	return &model.Movie{
		Kind:        model.ShowKindMovie,
		RecordingID: recordingDetails.RecordingID,
		Title:       recordingDetails.Title,
		RecordedOn:  recordingDetails.StartTime.Time,
		Description: recordingDetails.Description,
		MovieYear:   recordingDetails.MovieYear,
	}, nil
}

func NewSeries(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*model.Series, error) {
	if recordingDetails.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	return &model.Series{
		Kind:         model.ShowKindEpisode,
		CollectionID: recordingDetails.CollectionID,
		Title:        recordingDetails.Title,
		RecordedOn:   recordingDetails.StartTime.Time,
		Description:  collectionDetails.Description,
	}, nil
}

func NewEpisode(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*model.Episode, error) {
	if recordingDetails.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	var episodeNumber int
	if len(recordingDetails.EpisodeNum) > 0 {
		episodeNumber = recordingDetails.EpisodeNum[0]
	}
	return &model.Episode{
		Kind:               model.ShowKindEpisode,
		RecordingID:        recordingDetails.RecordingID,
		CollectionID:       collectionDetails.CollectionID,
		Title:              recordingDetails.Title,
		RecordedOn:         recordingDetails.StartTime.Time,
		Description:        collectionDetails.Description,
		OriginalAirDate:    recordingDetails.OriginalAirDate,
		SeasonNumber:       recordingDetails.SeasonNumber,
		EpisodeNumber:      episodeNumber,
		EpisodeTitle:       recordingDetails.Subtitle,
		EpisodeDescription: recordingDetails.Description,
	}, nil
}

func MergeEpisodes(shows []model.Show) []model.Show {

	combinedShowsMap := make(map[string]model.Show)

	for _, show := range shows {
		if show.GetKind() == model.ShowKindEpisode {
			episode := show.(*model.Episode)
			if existingSeries, exists := combinedShowsMap[episode.Title]; exists {
				series := existingSeries.(*model.Series)
				series.Episodes = append(series.Episodes, episode)
				if series.RecordedOn.Before(episode.RecordedOn) {
					series.RecordedOn = episode.RecordedOn
				}
			} else {
				combinedShowsMap[episode.Title] = &model.Series{
					Kind:         model.ShowKindSeries,
					CollectionID: episode.CollectionID,
					Title:        episode.Title,
					RecordedOn:   episode.RecordedOn,
					Description:  episode.Description,
					Episodes:     []*model.Episode{episode},
				}
			}
		} else if show.GetKind() == model.ShowKindMovie {
			combinedShowsMap[show.GetTitle()] = show
		}
	}

	combinedShows := make([]model.Show, 0, len(combinedShowsMap))
	for _, show := range combinedShowsMap {
		combinedShows = append(combinedShows, show)
	}
	sort.Slice(combinedShows, func(i, j int) bool {
		return combinedShows[i].GetTitle() < combinedShows[j].GetTitle()
	})

	return combinedShows
}
