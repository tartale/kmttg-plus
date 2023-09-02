package shows

import (
	"fmt"
	"sort"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

func New(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (model.Show, error) {

	switch recordingDetails.CollectionType {
	case message.CollectionTypeSeries:
		return newEpisode(recordingDetails, collectionDetails)

	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		return newMovie(recordingDetails)

	default:
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

}

func Cast(show model.Show) model.Show {
	switch show.GetKind() {
	case model.ShowKindMovie:
		return show.(*movie).Movie
	case model.ShowKindSeries:
		return show.(*series).Series
	case model.ShowKindEpisode:
		return show.(*episode).Episode
	}
	logz.Logger.Warn("unable to cast show to underlying type", zap.Any("kind", show.GetKind()))

	return show
}

func MergeEpisodes(shows []model.Show) []model.Show {

	combinedShowsMap := make(map[string]model.Show)

	for _, show := range shows {
		if show.GetKind() == model.ShowKindEpisode {
			if existingSeries, exists := combinedShowsMap[show.GetTitle()]; exists {
				episode := show.(*episode).Episode
				series := existingSeries.(*series)
				series.Episodes = append(series.Episodes, episode)
				if series.RecordedOn.Before(episode.RecordedOn) {
					series.RecordedOn = episode.RecordedOn
				}
			} else {
				episode := show.(*episode)
				combinedShowsMap[show.GetTitle()] = newSeries(episode)
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

type movie struct {
	*model.Movie
	recordingDetails *message.RecordingItem
}

func newMovie(recordingDetails *message.RecordingItem) (*movie, error) {
	if recordingDetails.CollectionType != message.CollectionTypeMovie && recordingDetails.CollectionType != message.CollectionTypeSpecial {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	return &movie{
		Movie: &model.Movie{
			ID:          recordingDetails.RecordingID,
			Kind:        model.ShowKindMovie,
			Title:       recordingDetails.Title,
			RecordedOn:  recordingDetails.StartTime.Time,
			Description: recordingDetails.Description,
			MovieYear:   recordingDetails.MovieYear,
		},
		recordingDetails: recordingDetails,
	}, nil
}

type series struct {
	*model.Series
	recordingDetails  *message.RecordingItem
	collectionDetails *message.CollectionItem
}

func newSeries(episode *episode) *series {

	return &series{
		Series: &model.Series{
			ID:          episode.SeriesID,
			Kind:        model.ShowKindSeries,
			Title:       episode.Title,
			RecordedOn:  episode.RecordedOn,
			Description: episode.Description,
			Episodes:    []*model.Episode{episode.Episode},
		},
		recordingDetails:  episode.recordingDetails,
		collectionDetails: episode.collectionDetails,
	}
}

type episode struct {
	*model.Episode
	recordingDetails  *message.RecordingItem
	collectionDetails *message.CollectionItem
}

func newEpisode(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*episode, error) {
	if recordingDetails.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	var episodeNumber int
	if len(recordingDetails.EpisodeNum) > 0 {
		episodeNumber = recordingDetails.EpisodeNum[0]
	}

	return &episode{
		Episode: &model.Episode{
			ID:                 recordingDetails.RecordingID,
			SeriesID:           recordingDetails.CollectionID,
			Kind:               model.ShowKindEpisode,
			Title:              recordingDetails.Title,
			RecordedOn:         recordingDetails.StartTime.Time,
			Description:        collectionDetails.Description,
			OriginalAirDate:    recordingDetails.OriginalAirDate,
			SeasonNumber:       recordingDetails.SeasonNumber,
			EpisodeNumber:      episodeNumber,
			EpisodeTitle:       recordingDetails.Subtitle,
			EpisodeDescription: recordingDetails.Description,
		},
		recordingDetails:  recordingDetails,
		collectionDetails: collectionDetails,
	}, nil
}
